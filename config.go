package dns_router

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"launchpad.net/goyaml"
)

type ApplicationConfig struct {
	DataDir     string
	BindAddr    string
	HealthCheck HealthCheckConfig

	// catch all backend
	Default DefaultBackend

	// backends based on pattern
	Backends []BackendConfig

	// simple override A records
	Hosts map[string]string

	// IP Alias maps an ip address to a shortname for logging convenience
	IPAlias map[string]string

	// logging config
	Logging LoggingConfig
}
type LoggingConfig struct {
	Enabled   bool
	Directory string
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
		Hosts:    make(map[string]string),
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

	err = c.Fixup()
	if err != nil {
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

func (c *ApplicationConfig) Fixup() error {

	if c.Logging.Directory == "" {
		c.Logging.Directory = filepath.Join(c.DataDir, "logs")
	}

	return nil

}
