package config

// Server defines the server configuration.
type Server struct {
	Addr          string
	Host          string
	Pprof         bool
	Root          string
	Cert          string
	Key           string
	StrictCurves  bool
	StrictCiphers bool
	Storage       string
}

// Metrics defines the metrics server configuration.
type Metrics struct {
	Addr  string
	Token string
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string
	Pretty bool
	Color  bool
}

// General defines the general configuration.
type General struct {
	Username string
	Password string
	Secret   string
}

// Config defines the general configuration.
type Config struct {
	Server  Server
	Metrics Metrics
	Logs    Logs
	General General
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
