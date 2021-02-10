package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-multierror"
	cons "github.com/myENA/consultant/v2"
	ndf "github.com/myENA/nodefflag"
	"github.com/nathanejohnson/conseption/putbackreader"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Prefix       string      // The prefix that we will be watching in the consul kv
	Orphan       bool        // whether we will be taking care of orphans or locals only
	ConsulPort   int         // port for consul - defaults to 8500
	ConsulConfig *api.Config // consul config
}

func NewDefaultConfig() *Config {
	return &Config{
		Prefix:     "/services",
		Orphan:     false,
		ConsulPort: 8500,
	}
}

var precomma = regexp.MustCompile(`^\s*,`)

func main() {

	flags := ndf.NewNDFlagSet(os.Args[0], flag.ExitOnError)
	conf := flags.NDString("config", "example.toml", "path to toml config file - optional")
	orphan := flags.NDBool("orphan", true, "orphan mode - default false")
	prefix := flags.NDString("prefix", "/services", "prefix to watch for service regs - default /services")
	cport := flags.NDInt("consulport", 8500, "tcp port for remote consul connections")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fatalf(err.Error())
	}

	cfg := NewDefaultConfig()
	if *conf != nil {
		tree, err := toml.LoadFile(**conf)
		if err != nil {
			fatalf("error loading config file %q: %v", **conf, err)
		}
		if err := tree.Unmarshal(cfg); err != nil {
			fatalf("error unmarshalling file %q: %v", **conf, err)
		}
	}
	if *orphan != nil {
		cfg.Orphan = **orphan
	}

	if *cport != nil {
		cfg.ConsulPort = **cport
	}

	if *prefix != nil {
		cfg.Prefix = **prefix
	}

	if cfg.ConsulConfig == nil {
		cfg.ConsulConfig = api.DefaultConfig()
	}

	cspt := &conseption{
		me:    os.Getenv("HOSTNAME"),
		cache: make(map[string][]byte),
		cfg:   cfg,
	}

	cspt.cc, err = cons.NewClient(cfg.ConsulConfig)

	if err != nil {
		fatalf("Error connecting to consul: %s\n", err)
	}

	cspt.node, err = cspt.cc.Agent().NodeName()
	if err != nil {
		fatalf("Cannot determine node name: %s", err)
	}

	wkp, err := cspt.cc.WatchKeyPrefix(cspt.cfg.Prefix, true, cspt.handler)
	if err != nil {
		fatalf("Error setting up watcher: %s\n", err)
	}

	svcs, _, err := cspt.cc.Catalog().Services(&api.QueryOptions{AllowStale: true})
	if err != nil {
		fatalf("Error querying catalog: %s\n", err)
	}

	// Go through all services in the catalog, and deregister anything that
	// should go to the local host agent.
	for k := range svcs {
		css, _, err := cspt.cc.Catalog().Service(k, "", &api.QueryOptions{AllowStale: true})
		if err != nil {
			fmt.Printf("got error querying catalog: %s\n", err)
		}
		for _, cs := range css {
			if cs.ServiceAddress == cspt.me && cs.Node != cspt.node {
				err = cspt.deregRemote(cs)
				if err != nil {
					fmt.Printf("got error derigstering: %s\n", err)
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

type conseption struct {
	sync.Mutex
	cc    *cons.Client
	me    string
	node  string
	cache map[string][]byte
	cfg   *Config
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
	var regs []*api.AgentServiceRegistration

	// ensure we have kvpairs
	kvps, ok := raw.(api.KVPairs)
	if !ok {

		fmt.Println("not KVPairs!")
		return
	}
	news := make(map[string]bool)

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
			if svc.Address == cspt.me {
				h := md5.Sum(kvp.Value)
				if ch, ok := cspt.cache[svc.ID]; ok {
					if bytes.Equal(ch, h[:]) {
						// no change
						news[svc.ID] = false
						continue
					}
				}
				news[svc.ID] = true
				regs = append(regs, svc)
				cspt.cache[svc.ID] = h[:]
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

func parseServiceRegs(val []byte) ([]*api.AgentServiceRegistration, error) {
	// by the end of his func, this will either be nil or an *multierr.Error
	var sumOfAllErrs error

	// Try services struct
	ss := &services{}

	if err := json.Unmarshal(val, ss); err == nil {
		return ss.Services, nil
	}

	// now try a list
	if err := json.Unmarshal(val, &ss.Services); err == nil {
		return ss.Services, nil
	}

	// now try serial json objects.
	pbr := putbackreader.NewPutBackReader(bytes.NewReader(val))
	jd := json.NewDecoder(pbr)
	buf := new(bytes.Buffer)

	for {
		asr := &api.AgentServiceRegistration{}
		err := jd.Decode(asr)

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
				sumOfAllErrs = multierror.Append(sumOfAllErrs, fmt.Errorf("bad read: %s", string(b)))
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
					sumOfAllErrs = multierror.Append(sumOfAllErrs, fmt.Errorf("got final error: %w", err))
				}
				break
			}
		}
		ss.Services = append(ss.Services, asr)

		if !jd.More() {
			break
		}
	}

	return ss.Services, sumOfAllErrs
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
	fmt.Printf(format, args...)
	os.Exit(1)
}
