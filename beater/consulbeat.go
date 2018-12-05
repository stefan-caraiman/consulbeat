package beater

import (
	"fmt"
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
	agent :=  consulClient.Agent()


	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
		// Get Node and Service checks
		info, _ := agent.Self()
		// Use Catalog
		name := info["Config"]["NodeName"].(string)
		checks, meta, _ := health.Node(name, nil)
		fmt.Println(checks[0].Node, checks[0].Name, checks[0].Status, checks[0].Output)
		//logp.Info(string(checks))
		logp.Info(string(meta.LastIndex))
		// Parse checks
		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": 1337,
			},
		}
		bt.client.Publish(event)
		logp.Info("Event sent")
	}
}

// Stop stops consulbeat.
func (bt *Consulbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
