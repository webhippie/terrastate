package config

// Server defines the server configuration.
type Server struct {
	Addr          string `mapstructure:"addr"`
	Host          string `mapstructure:"host"`
	Pprof         bool   `mapstructure:"pprof"`
	Root          string `mapstructure:"root"`
	Cert          string `mapstructure:"cert"`
	Key           string `mapstructure:"key"`
	StrictCurves  bool   `mapstructure:"strict_curves"`
	StrictCiphers bool   `mapstructure:"strict_ciphers"`
	Storage       string `mapstructure:"storage"`
}

// Metrics defines the metrics server configuration.
type Metrics struct {
	Addr  string `mapstructure:"addr"`
	Token string `mapstructure:"token"`
}

// Logs defines the level and color for log configuration.
type Logs struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
}

// Encryption defines the encryption configuration.
type Encryption struct {
	Secret string `mapstructure:"secret"`
}

// Access defines the access configuration.
type Access struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Config defines the general configuration.
type Config struct {
	Server     Server     `mapstructure:"server"`
	Metrics    Metrics    `mapstructure:"metrics"`
	Logs       Logs       `mapstructure:"log"`
	Encryption Encryption `mapstructure:"encryption"`
	Access     Access     `mapstructure:"access"`
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
