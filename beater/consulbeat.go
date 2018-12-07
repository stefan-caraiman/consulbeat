package beater

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/stefan-caraiman/consulbeat/config"
)

// Consulbeat configuration.
type Consulbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of consulbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Consulbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts consulbeat.
func (bt *Consulbeat) Run(b *beat.Beat) error {
	logp.Info("Consulbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	// Initialize Consul client
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	//health := consulClient.Health()
	//agent :=  consulClient.Agent()
	catalog := consulClient.Catalog()

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		// Get Node and Service checks
		// Use Catalog
		// nodes, meta, err := catalog.Nodes(nil)
		// if err != nil {
		// 	panic(err)
		// }

		// if meta.LastIndex == 0 {
		// 	fmt.Printf("Bad: %v", meta)
		// }

		// if len(nodes) == 0 {
		// 	fmt.Printf("Bad: %v", nodes)
		// }
		// fmt.Println("Nodes: ", nodes)
		// for key, value := range nodes {
		// 	fmt.Printf("Key[%s] value [%s] \n", key, value)
		// }
		// fmt.Println(nodes[0].Datacenter)
		// fmt.Println(nodes[0].Address)
		// services
		services, meta_services, err := catalog.Services(nil)
		if err != nil {
			panic(err)
		}

		if meta_services.LastIndex == 0 {
			fmt.Printf("Bad: %v", meta_services)
		}

		if len(services) == 0 {
			fmt.Printf("Bad: %v", services)
		}

		var events = []beat.Event{}
		for service, _ := range services {
			//catalog_check, meta, err := health.Service(service, "", false, nil)
			checks, meta, err := catalog.Service(service, "", nil)
			if err != nil {
				panic(err)
			}

			if meta.LastIndex == 0 {
				fmt.Printf("Bad: %v", meta)
			}

			if len(checks) == 0 {
				fmt.Printf("Bad: %v", checks)
			}

			fmt.Println("Service name is: ", service)
			fmt.Println("Node status: ", checks[0].Node)
			for _, check_value := range checks {
				fmt.Println("Checked value is ", check_value)
				fmt.Println(check_value.ServiceAddress)
				fmt.Println(check_value.Checks)
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"node":         check_value.Node,
						"Datacenter":   check_value.Datacenter,
						"AgentAddress": check_value.Address,
						// "CheckID":        health_check_def.CheckID,
						// "CheckName":      health_check_def.Name,
						// "Status":         health_check_def.Status,
						// "Output":         health_check_def.Output,
						"ServiceID":      check_value.ServiceID,
						"ServicePort":    check_value.ServicePort,
						"ServiceAddress": check_value.ServiceAddress,
						"ServiceName":    check_value.ServiceName,
						"ServiceTags":    strings.Join(check_value.ServiceTags, ","),
					},
				}
				events = append(events, event)
			}

		}
		bt.client.PublishAll(events)
		logp.Info("Events sent")
	}
}

// Stop stops consulbeat.
func (bt *Consulbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
