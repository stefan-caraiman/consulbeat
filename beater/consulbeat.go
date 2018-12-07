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
	catalog := consulClient.Catalog()

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		services, meta_services, err := catalog.Services(nil)
		if err != nil {
			panic(err)
		}

		if meta_services.LastIndex == 0 {
			logp.Info("Bad: %v", meta_services)
		}

		if len(services) == 0 {
			logp.Info("Bad: %v", services)
		}

		var events = []beat.Event{}
		for service, _ := range services {
			//catalog_check, meta, err := health.Service(service, "", false, nil)
			checks, meta, err := catalog.Service(service, "", nil)
			if err != nil {
				panic(err)
			}

			if meta.LastIndex == 0 {
				logp.Info("Bad: %v", meta)
			}

			if len(checks) == 0 {
				logp.Info("Bad: %v", checks)
			}

			for _, service_check := range checks {
				health_service_check, _, _ := health.Checks(service_check.ServiceName, nil)
				for _, check := range health_service_check {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields: common.MapStr{
							"node":            service_check.Node,
							"datacenter":      service_check.Datacenter,
							"agent.address":   service_check.Address,
							"check.id":        check.CheckID,
							"check.name":      check.Name,
							"check.status":    check.Status,
							"check.output":    check.Output,
							"service.name":      check.ServiceName,
							"service.tags":    strings.Join(check.ServiceTags, ","),
							"service.id":      service_check.ServiceID,
							"service.port":    service_check.ServicePort,
							"service.address": service_check.ServiceAddress,
						},
					}
					events = append(events, event)
				}
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
