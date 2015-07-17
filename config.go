package dns_router

import (
	"bytes"
	"fmt"
	"launchpad.net/goyaml"
	"os"
)

type ApplicationConfig struct {
	BindAddr    string
	HealthCheck HealthCheckConfig
	Default     DefaultBackend
	Backends    []BackendConfig
	Hosts       map[string][]string
}
type HealthCheckConfig struct {
	Interval int
}
type DefaultBackend struct {
	Servers []string
}

type BackendConfig struct {
	Pattern     string
	HealthCheck bool
	Servers     []string
}

func Default() *ApplicationConfig {
	c := &ApplicationConfig{
		BindAddr: "127.0.0.1:53",
	}
	return c
}
func (c *ApplicationConfig) LoadYaml(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)
	_, err = b.ReadFrom(f)
	if err != nil {
		return err
	}

	if err := c.LoadYamlBuffer(b.Bytes()); err != nil {
		return err
	}

	return nil
}

func (c *ApplicationConfig) LoadYamlBuffer(buf []byte) error {
	err := goyaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApplicationConfig) PrintConfig() {
	d, err := goyaml.Marshal(c)
	if err != nil {
		fmt.Println("Marshal error", err)
		return
	}
	fmt.Println(string(d))
}
