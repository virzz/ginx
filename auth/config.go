package auth

import (
	"fmt"

	"github.com/spf13/pflag"
)

func FlagSet(name string) *pflag.FlagSet {
	fs := pflag.NewFlagSet("http.auth", pflag.ContinueOnError)

	fs.Bool("http.auth.enabled", false, "Enable Auth")

	fs.Int("http.auth.maxage", 7*24*3600, "HTTP Auth MaxAge")
	fs.String("http.auth.secret", "", "Session Store Secret")

	fs.String("http.auth.host", "127.0.0.1", "Redis Store Host")
	fs.Int("http.auth.port", 6379, "Redis Store Port")
	fs.Int("http.auth.db", 0, "Redis Store DB")
	fs.String("http.auth.pass", "", "Redis Store Password")

	return fs
}

//go:generate structx -struct Config
type Config struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
	MaxAge  int  `json:"maxage" yaml:"maxage"` // Expire(MaxAge)

	// Cookie
	Secret string `json:"secret" yaml:"secret"` // Cookie Secret

	// Token | Redis
	Host string `json:"host" yaml:"host"` // Redis Address
	Port int    `json:"port" yaml:"port"` // Redis Port
	Pass string `json:"pass" yaml:"pass"` // Redis Password
	DB   int    `json:"db" yaml:"db"`     // Redis DB
}

func (s *Config) WithAddr(v string) *Config {
	s.Host = v
	return s
}

func (s *Config) WithPort(v int) *Config {
	s.Port = v
	return s
}

func (s *Config) WithPass(v string) *Config {
	s.Pass = v
	return s
}

func (s *Config) WithDB(v int) *Config {
	s.DB = v
	return s
}

func (s *Config) WithMaxAge(v int) *Config {
	s.MaxAge = v
	return s
}

func (s *Config) Addr() string { return fmt.Sprintf("%s:%d", s.Host, s.Port) }
