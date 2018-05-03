package config

// Server defines the server configuration.
type Server struct {
	Private       string
	Public        string
	Host          string
	Root          string
	Cert          string
	Key           string
	StrictCurves  bool
	StrictCiphers bool
	Templates     string
	Assets        string
	Storage       string
}

// Logs defines the logging configuration.
type Logs struct {
	Level   string
	Colored bool
	Pretty  bool
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
	Logs    Logs
	General General
}

// New prepares a new default configuration.
func New() *Config {
	return &Config{}
}
