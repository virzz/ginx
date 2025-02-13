// Code generated by github.com/virzz/structx. DO NOT EDIT.

package ginx

func (s *Config) WithSystem(v string)     { s.System = v }
func (s *Config) WithEndpoint(v string)   { s.Endpoint = v }
func (s *Config) WithAddr(v string)       { s.Addr = v }
func (s *Config) WithPort(v int)          { s.Port = v }
func (s *Config) WithHeaders(v []string)  { s.Headers = v }
func (s *Config) WithPprof(v bool)        { s.Pprof = v }
func (s *Config) WithMetrics(v bool)      { s.Metrics = v }
func (s *Config) WithOrigins(v []string)  { s.Origins = v }
func (s *Config) WithDebug(v bool)        { s.Debug = v }
func (s *Config) WithRequestID(v bool)    { s.RequestID = v }
func (s *Config) WithStore(v StoreConfig) { s.Store = v }
