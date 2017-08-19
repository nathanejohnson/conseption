package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/hashicorp/consul/api"
	cons "github.com/myENA/consultant"
)

func main() {

	cspt := conseption{
		conf: api.DefaultConfig(),
		me:   os.Getenv("HOSTNAME"),
	}
	var err error

	cspt.cc, err = cons.NewClient(cspt.conf)

	if err != nil {
		Fatalf("Error connecting to consul: %s\n", err)
		os.Exit(1)
	}

	wkp, err := cspt.cc.WatchKeyPrefix("/services", true, cspt.handler)
	if err != nil {
		Fatalf("Error setting up watcher: %s\n", err)
	}

	wkp.Run(cspt.conf.Address)

}

type conseption struct {
	conf *api.Config
	cc   *cons.Client
	me   string
}

type services struct {
	Services []*api.AgentServiceRegistration
}

func (cspt *conseption) handler(idx uint64, raw interface{}) {
	v := reflect.ValueOf(raw)
	if v.Kind() == reflect.Slice {
		t := v.Type()
		if t.Name() == "KVPairs" {
			for i := 0; i < v.Len(); i++ {
				kvp := (*api.KVPair)(unsafe.Pointer(v.Index(i).Pointer()))
				ss := &services{}
				err := json.Unmarshal(kvp.Value, ss)
				if err != nil {
					// try doing it one at a time
					err = nil
					jd := json.NewDecoder(bytes.NewReader(kvp.Value))
					for jd.More() {
						s := &api.AgentServiceRegistration{}
						err = jd.Decode(s)
						if err != nil {
							continue
						}
						ss.Services = append(ss.Services, s)
					}
				}

				_ = cspt.deregisterAllServices() // todo log this
				for _, svc := range ss.Services {
					if svc.Address == cspt.me {
						fmt.Printf("I'd totally register %#v\n", svc)
						// cspt.cc.Agent().ServiceRegister(svc)
					}
				}

			}
		}
	} else {
		fmt.Printf("not pointer: %#v\n", raw)
		fmt.Printf("kind: %s\n", v.Kind().String())
	}
}

func (cspt *conseption) deregisterAllServices() error {
	a := cspt.cc.Agent()
	services, err := a.Services()
	if err != nil {
		return err
	}
	var errs []string
	for _, s := range services {
		// err = a.ServiceDeregister(s.ID)
		fmt.Printf("I'd totally derigster %s\n", s.ID)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if errs != nil {
		return fmt.Errorf("errors deregistering service: %s", strings.Join(errs, ","))
	}
	return nil
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args)
	os.Exit(1)
}
