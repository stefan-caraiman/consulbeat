// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	ConsulAddress string `config:"consuladdress"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	ConsulAddress: "127.0.0.1:8500",

}
