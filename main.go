package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	cons "github.com/myENA/consultant"
	ndf "github.com/myENA/nodefflag"
	"github.com/nathanejohnson/conseption/putbackreader"
)

type Config struct {
	// The prefix that we will be watching in the consul kv.  This is where our agent
	// service registrations live.
	Prefix string

	// whether we will be taking care of orphans or locals only.  Defaults to false,
	// meaning locals only
	Orphanage bool

	// port for remote consul agents - defaults to 8500 (why isn't this in the catalog?)
	ConsulPort int

	// ConsulConfig - consul api.Config struct.  if nil, a sensible default will be used.
	ConsulConfig *api.Config // consul config

	// TagPrefix for service registration
	TagPrefix string

	// ServiceName for consul registration
	ServiceName string

	// TTL for service registration
	TTLInterval api.ReadableDuration

	// if true, deregister our conseption service from local agent.  defaults to true.
	DeregOnExit bool `toml",omitempty"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Prefix:       "/services",
		Orphanage:    false,
		ConsulPort:   8500,
		ConsulConfig: api.DefaultConfig(),
		TagPrefix:    "conseption",
		ServiceName:  "conseption",
		TTLInterval:  api.ReadableDuration(time.Second * 30),
		DeregOnExit:  true,
	}
}

var precomma = regexp.MustCompile(`^\s*,`)

func main() {

	flags := ndf.NewNDFlagSet(os.Args[0], flag.ExitOnError)
	conf := flags.ZVString("config", "-config=/path/to/some/config", "path to toml config file - optional")
	orphan := flags.ZVBool("orphan", true, "orphan mode - default false")
	prefix := flags.ZVString("prefix", "-prefix=/services", "prefix to watch for service regs - default /services")
	cport := flags.NDInt("consulport", 8500, "tcp port for remote consul connections")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fatalf(err.Error())
	}

	cfg := NewDefaultConfig()
	if *conf != "" {
		_, err = toml.DecodeFile(*conf, cfg)
		if err != nil {
			fatalf(err.Error())
		}
	}

	cfg.Orphanage = *orphan

	if *cport != nil {
		cfg.ConsulPort = **cport
	}

	if *prefix != "" {
		cfg.Prefix = *prefix
	}

	if cfg.ConsulConfig == nil {
		cfg.ConsulConfig = api.DefaultConfig()
	}

	ips, err := ifaces()
	if err != nil {
		fmt.Println("Error fetching interfaces:", err)
	}

	cspt := &conseption{
		me:       strings.ToLower(os.Getenv("HOSTNAME")),
		cache:    make(map[string]*cacheEntry),
		cfg:      cfg,
		localIPs: ips,
	}

	cspt.cc, err = cons.NewClient(cfg.ConsulConfig)

	if err != nil {
		fatalf("Error connecting to consul: %s\n", err)
	}

	// register our own service handler.

	asr := &api.AgentServiceRegistration{
		ID:   cspt.cfg.ServiceName,
		Name: cspt.cfg.ServiceName,
	}

	if cspt.cfg.Orphanage {
		asr.Tags = []string{"orphanage"}
	} else {
		asr.Tags = []string{"localagent"}
	}

	err = cspt.cc.Agent().ServiceRegister(asr)
	if err != nil {
		fatalf("Could not register our service: %s\n", err)
	}
	cspt.chkid = cspt.cfg.ServiceName + "_ttl"
	err = cspt.cc.Agent().CheckRegister(&api.AgentCheckRegistration{
		ID:   cspt.chkid,
		Name: cspt.chkid,
		AgentServiceCheck: api.AgentServiceCheck{
			TTL: cspt.cfg.TTLInterval.String(),
		},
	})

	go func() {
		sleepint := time.Duration(cspt.cfg.TTLInterval / 2)
		for {
			err := cspt.cc.Agent().UpdateTTL(cspt.chkid, "saul goodman", "passing")
			if err != nil {
				fmt.Println("Got error updating TTL", err)
			}
			time.Sleep(sleepint)
		}
	}()

	cspt.node, err = cspt.cc.Agent().NodeName()
	if err != nil {
		fatalf("Cannot determine node name: %s", err)
	}

	// Seed the cache

	kvps, _, err := cspt.cc.KV().List(cspt.cfg.Prefix, &api.QueryOptions{AllowStale: true})

	if err != nil {
		fatalf("Error fetching from prefix: %s\n", err)
	}

	cspt.handler(0, kvps)

	wkp, err := cspt.cc.WatchKeyPrefix(cspt.cfg.Prefix, true, cspt.handler)
	if err != nil {
		fatalf("Error setting up watcher: %s\n", err)
	}

	svcs, _, err := cspt.cc.Catalog().Services(&api.QueryOptions{AllowStale: true})
	if err != nil {
		fatalf("Error querying catalog: %s\n", err)
	}

	// TODO upgrade this to handle situations where a service is registered with an IP that maps to our address
	// or an address that maps to our IP
	if !cspt.cfg.Orphanage {
		for k := range svcs {
			css, _, err := cspt.cc.Catalog().Service(k, "", &api.QueryOptions{AllowStale: true})
			if err != nil {
				fmt.Printf("got error querying catalog: %s\n", err)
			}
			for _, cs := range css {
				if entry, ok := cspt.cache[cskey(cs)]; ok && cspt.isItI(cs.Address) && cs.Node != cspt.node {
					_, _, err = cspt.cc.Event().Fire(&api.UserEvent{
						Name:    cspt.cfg.Prefix + "_takeover",
						Payload: entry.sum,
					}, nil)
					if err != nil {
						fmt.Printf("got error firing takeover event: %s\n", err)
					}
					err = cspt.deregRemote(cs)
					if err != nil {
						fmt.Printf("got error derigstering: %s\n", err)
					}
				}
			}
		}
	}

	// Start the runner, which will get an initial full kv dump of everything
	// under the kv prefix.
	err = wkp.Run(cspt.cfg.ConsulConfig.Address)
	if err != nil {
		fatalf("%s\n", err)
	}
}

func (cspt *conseption) isItI(addr string) bool {
	// first check to see if addr is an IP
	var (
		ips []net.IP
		err error
	)

	if strings.ToLower(addr) == cspt.me {
		return true
	}

	ip := net.ParseIP(addr)
	if ip == nil {
		ips, err = net.LookupIP(addr)
		if err != nil {
			return false
		}
	} else {
		ips = append(ips, ip)
	}
	// TODO - make this a map lookup on localIPs
	for _, lip := range cspt.localIPs {
		for _, sip := range ips {
			if lip.Equal(sip) {
				return true
			}
		}
	}
	return false

}

func ifaces() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ips = append(ips, v.IP)
			case *net.IPAddr:
				ips = append(ips, v.IP)
			default:
				return nil, fmt.Errorf("unexpected address from interface %s: %s", iface.Name, addr.String())
			}
		}
	}
	return ips, nil
}

type conseption struct {
	sync.Mutex
	cc       *cons.Client
	me       string
	node     string
	cache    map[string]*cacheEntry
	cfg      *Config
	localIPs []net.IP
	chkid    string
}

type cacheEntry struct {
	sum []byte
	asr *api.AgentServiceRegistration
}

type services struct {
	Services []*api.AgentServiceRegistration
}

func (cspt *conseption) deregRemote(se *api.CatalogService) error {
	// shallow copy of conf
	conf := &api.Config{}
	*conf = *cspt.cfg.ConsulConfig

	// get info on the node
	cn, _, err := cspt.cc.Catalog().Node(se.Node, &api.QueryOptions{AllowStale: true})
	if err != nil {
		return err
	}

	conf.Address = fmt.Sprintf("%s:%d", cn.Node.Address, cspt.cfg.ConsulPort)

	rcc, err := cons.NewClient(conf)
	if err != nil {
		return err
	}
	fmt.Printf("deregistering %s from %s\n", se.ServiceName, se.Node)

	return rcc.Agent().ServiceDeregister(se.ServiceID)
}

func (cspt *conseption) handler(idx uint64, raw interface{}) {

	cspt.Lock()
	defer cspt.Unlock()
	kvps, ok := raw.(api.KVPairs)
	var regs []*api.AgentServiceRegistration
	if !ok {
		fmt.Println("not KVPairs!")
		return
	}
	news := make(map[string]bool)

	qo := &api.QueryOptions{AllowStale: true}

	for _, kvp := range kvps {
		fmt.Printf("handling %s\n", kvp.Key)
		svcs, err := parseServiceRegs(kvp.Value)
		if err != nil {
			fmt.Printf("error parsing service reg: %s\n", err)
			if svcs == nil {
				return
			}
		}
		for _, svc := range svcs {
			k := askey(svc)

			//TODO - deregister if it's mine and
			// someone elsehas it.
			if cspt.isItI(svc.Address) {
				h := md5.Sum(kvp.Value)
				if ch, ok := cspt.cache[k]; ok {
					if bytes.Equal(ch.sum, h[:]) {
						// no change
						news[k] = false
						continue
					}
				}
				news[k] = true
				regs = append(regs, svc)
				cspt.cache[k] = &cacheEntry{
					sum: h[:],
					asr: svc,
				}
			} else if cspt.cfg.Orphanage {
				var tag string
				if len(svc.Tags) > 0 {
					tag = svc.Tags[0]
				}
				sents, _, err := cspt.cc.Health().Service(svc.Name, tag, false, qo)
				if err != nil {
					fmt.Println("Error fetching services", err)
					continue
				}
				for _, sent := range matchesTags(sents, svc.Tags) {

				}
			}
		}
	}
	var deregs []string
	for k := range cspt.cache {
		if u, ok := news[k]; ok {
			if u {
				// we were updated, so deregister and re-register
				deregs = append(deregs, k)
			} // else no change, no op
		} else {
			// key deleted
			deregs = append(deregs, k)
			delete(cspt.cache, k)
		}
	}

	for _, dr := range deregs {
		fmt.Printf("deregistering %s\n", dr)
		err := cspt.cc.Agent().ServiceDeregister(dr)
		if err != nil {
			fmt.Printf("got error deregistering %s: %s\n", dr, err)
		}
	}

	for _, svc := range regs {
		var err error
		fmt.Printf("registering %s\n", svc.ID)
		err = cspt.cc.Agent().ServiceRegister(svc)
		if err != nil {
			fmt.Printf("error returned from registering service: %s\n", err)
		}
	}
}

// prune service entries so that the only entries returned match *all* the tags passed in.
func matchesTags(ses []*api.ServiceEntry, tags []string) []*api.ServiceEntry {
	var rv []*api.ServiceEntry

	tagmap := make(map[string]bool)

	for _, t := range tags {
		tagmap[t] = true
	}
OUTER:
	for _, se := range ses {
		if len(tagmap) != len(se.Service.Tags) {
			continue
		}
		for _, t := range se.Service.Tags {
			if !tagmap[t] {
				continue OUTER
			}
		}
		rv = append(rv, se)
	}

	return rv

}

func askey(svc *api.AgentServiceRegistration) string {
	id := svc.ID
	if svc.ID == "" {
		id = svc.Name
	}
	return fmt.Sprintf("%s;%s:%d", id, svc.Address, svc.Port)
}

func cskey(cs *api.CatalogService) string {
	return fmt.Sprintf("%s;%s:%d", cs.ServiceID, cs.Address, cs.ServicePort)
}

func parseServiceRegs(val []byte) ([]*api.AgentServiceRegistration, error) {
	var errors []string
	var err error
	// Try services struct
	ss := &services{}
	err = json.Unmarshal(val, ss)
	if err == nil {
		return ss.Services, nil
	}

	// now try a list
	err = json.Unmarshal(val, &ss.Services)
	if err == nil {
		return ss.Services, nil
	}

	// now try serial json objects.
	pbr := putbackreader.NewPutBackReader(bytes.NewReader(val))
	jd := json.NewDecoder(pbr)
	buf := new(bytes.Buffer)

	for {
		asr := &api.AgentServiceRegistration{}
		err = jd.Decode(asr)

		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			// Handle the case where we have comma separated json
			// objects.
			buf.Reset()
			_, _ = buf.ReadFrom(jd.Buffered())
			b := buf.Bytes()
			m := precomma.FindIndex(b)
			if m == nil {
				errors = append(errors, fmt.Sprintf("bad read: %s\n", string(b)))
				break
			}

			// Take the comma off, put the already-read parts of the stream
			// back, and make a new decoder.  All this work to subtract
			// a fucking wayward comma from the stream.
			pbr.SetBackBytes(b[m[1]:])
			jd = json.NewDecoder(pbr)

			err = jd.Decode(asr)
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					errors = append(errors, fmt.Sprintf("got final error: %s\n", err))
				}
				break
			}
		}
		ss.Services = append(ss.Services, asr)

		if !jd.More() {
			break
		}
	}

	if len(errors) > 0 {
		err = fmt.Errorf("Errors: %s", strings.Join(errors, ","))
	}
	return ss.Services, err
}

func (cspt *conseption) deregisterAllLocalServices() error {
	a := cspt.cc.Agent()
	services, err := a.Services()
	if err != nil {
		return err
	}
	var errs []string
	for _, s := range services {
		err = a.ServiceDeregister(s.ID)
		fmt.Printf("deregistering %s\n", s.ID)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if errs != nil {
		return fmt.Errorf("errors deregistering service: %s", strings.Join(errs, ","))
	}
	return nil
}

func fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args)
	os.Exit(1)
}
