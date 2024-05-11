package config

const (
	defaultPort = "4221"
	defaultDir  = "/"
)

type Config struct {
	port string
	dir  string
}

type Option func(*Config)

func NewConfig(options ...Option) *Config {
	cfg := &Config{
		port: defaultPort,
		dir:  defaultDir,
	}

	for _, option := range options {
		option(cfg)
	}

	return cfg
}

func (cfg *Config) Port() string {
	return cfg.port
}

func (cfg *Config) Dir() string {
	return cfg.dir
}

func WithPort(port string) Option {
	return func(c *Config) {
		c.port = port
	}
}

func WithDir(dir string) Option {
	return func(c *Config) {
		c.dir = dir
	}
}
