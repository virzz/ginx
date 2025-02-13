package ginx

import "github.com/spf13/pflag"

func FlagSet(defaultPort int) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("http", pflag.ContinueOnError)
	flagSet.String("http.endpoint", "", "HTTP Domain Endpoint")
	flagSet.String("http.addr", "127.0.0.1", "HTTP Listen Address")
	flagSet.Int("http.port", defaultPort, "HTTP Listen Port")

	flagSet.StringSlice("http.origins", []string{"*"}, "HTTP CORS: Allow Origins")
	flagSet.StringSlice("http.headers", []string{"Authorization"}, "HTTP CORS: Allow Headers")
	flagSet.Bool("http.debug", false, "HTTP Debug Mode")

	flagSet.Bool("http.captcha", false, "Enable Any Captaha")
	flagSet.Bool("http.pprof", false, "Enable PProf")
	flagSet.Bool("http.requestid", false, "Enable HTTP RequestID")
	flagSet.Bool("http.metrics", false, "Enable Metrics")

	flagSet.Bool("http.store.enabled", false, "Enable Redis Store")
	flagSet.String("http.store.addr", "127.0.0.1", "HTTP Session Store(Redis) Address")
	flagSet.Int("http.store.port", 6379, "HTTP Session Store(Redis) Port")
	flagSet.Int("http.store.db", 7, "HTTP Session Store(Redis) DB")
	flagSet.String("http.store.pass", "", "HTTP Session Store(Redis) Password")

	flagSet.String("http.upgrade", "", "HTTP Upgrade Token")
	return flagSet
}

//go:generate structx -struct StoreConfig
type StoreConfig struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Addr    string `json:"addr" yaml:"addr"`
	Port    int    `json:"port" yaml:"port"`
	Pass    string `json:"pass" yaml:"pass"`
	DB      int    `json:"db" yaml:"db"`
}

//go:generate structx -struct Config
type Config struct {
	Endpoint  string      `json:"endpoint" yaml:"endpoint"`
	Addr      string      `json:"addr" yaml:"addr"`
	Port      int         `json:"port" yaml:"port"`
	Origins   []string    `json:"origins" yaml:"origins"`
	Headers   []string    `json:"headers" yaml:"headers"`
	Debug     bool        `json:"debug" yaml:"debug"`
	Captcha   bool        `json:"captcha" yaml:"captcha"`
	Pprof     bool        `json:"pprof" yaml:"pprof"`
	RequestID bool        `json:"requestid" yaml:"requestid"`
	Metrics   bool        `json:"metrics" yaml:"metrics"`
	Store     StoreConfig `json:"store" yaml:"store"`
	System    string      `json:"system" yaml:"system"`
}

var Conf *Config
