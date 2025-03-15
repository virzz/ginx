package ginx

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/virzz/ginx/auth"
)

func FlagSet(defaultPort int) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("http", pflag.ContinueOnError)
	flagSet.String("http.prefix", "/", "HTTP API Route Prefix")

	flagSet.String("http.endpoint", "", "HTTP Domain Endpoint")
	flagSet.String("http.addr", "127.0.0.1", "HTTP Listen Address")
	flagSet.Int("http.port", defaultPort, "HTTP Listen Port")

	flagSet.StringSlice("http.origins", []string{"*"}, "HTTP CORS: Allow Origins")
	flagSet.StringSlice("http.headers", []string{"Authorization"}, "HTTP CORS: Allow Headers")
	flagSet.Bool("http.debug", false, "HTTP Debug Mode")

	flagSet.Bool("http.pprof", false, "Enable PProf")
	flagSet.Bool("http.requestid", false, "Enable HTTP RequestID")
	flagSet.Bool("http.metrics", false, "Enable Metrics")

	flagSet.AddFlagSet(auth.FlagSet("http.auth"))

	flagSet.String("http.system", "", "HTTP System Token")
	return flagSet
}

//go:generate structx -struct Config
type Config struct {
	version string `json:"-" yaml:"-"`
	commit  string `json:"-" yaml:"-"`

	System    string      `json:"system" yaml:"system"`
	Prefix    string      `json:"prefix" yaml:"prefix"`
	Endpoint  string      `json:"endpoint" yaml:"endpoint"`
	Addr      string      `json:"addr" yaml:"addr"`
	Port      int         `json:"port" yaml:"port"`
	Origins   []string    `json:"origins" yaml:"origins"`
	Headers   []string    `json:"headers" yaml:"headers"`
	Debug     bool        `json:"debug" yaml:"debug"`
	Pprof     bool        `json:"pprof" yaml:"pprof"`
	RequestID bool        `json:"requestid" yaml:"requestid"`
	Metrics   bool        `json:"metrics" yaml:"metrics"`
	Auth      auth.Config `json:"auth" yaml:"auth"`
}

func (c *Config) GetEndpoint() string {
	if c.Endpoint != "" {
		c.Endpoint = c.GetAddr()
	}
	return c.Endpoint
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
