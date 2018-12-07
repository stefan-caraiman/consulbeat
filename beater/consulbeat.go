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
	health := consulClient.Health()
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
			checks, meta, err := health.Service(service, "", false, nil)
			catalog_check, _, _ := catalog.Service(service, "", nil)
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
				for _, health_check_def := range check_value.Checks {
					// Marshal back to a map
					event := beat.Event{
						Timestamp: time.Now(),
						Fields: common.MapStr{
							"type": "service", // REFACTOR
							"ok": catalog_check[0].ServiceAddress,
							"node": check_value.Node.Node,
							"Datacenter": check_value.Node.Datacenter,
							"AgentAddress": check_value.Node.Address,
							"Port": check_value.Service.Port,
							"Service": check_value.Service.Service,
							"CheckID": health_check_def.CheckID,
							"CheckName": health_check_def.Name,
							"Status": health_check_def.Status,
							"Output": health_check_def.Output,
							"ServiceID": health_check_def.ServiceID,
							"ServiceName": health_check_def.ServiceName,
							"ServiceTags": strings.Join(health_check_def.ServiceTags, ","),
						},
					}
					events = append(events, event)
				}
			}
		}

		// Testing node checks
		//info, _ := agent.Self()
		//name := info["Config"]["NodeName"].(string)
		//checks, meta, _ := health.Node(name, nil)
		// fmt.Println("Node checks: ", checks[0].Node)
		// fmt.Println("Name: ", checks[0].Name)
		// fmt.Println("Status check: ", checks[0].Status)
		// fmt.Println("Output: ", checks[0].Output)
		//logp.Info(string(checks))
		//logp.Info(string(meta.LastIndex))
		// Parse checks
		bt.client.PublishAll(events)
		logp.Info("Events sent")
	}
}

// Stop stops consulbeat.
func (bt *Consulbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
