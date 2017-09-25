package config

type server struct {
	Host       string
	Addr       string
	Cert       string
	Key        string
	Storage    string
	Pprof      bool
	Prometheus bool
}

type general struct {
	Username string
	Password string
}

var (
	// LogLevel defines the log level used by our logging package.
	LogLevel string

	// Server represents the information about the server bindings.
	Server = &server{}

	// General represents the information about the general bindings.
	General = &general{}
)
