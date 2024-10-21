package conf

type Config struct {
	GCTrigger func() bool
	EnableTCO bool
}

func New(opts ...Option) *Config {
	c := &Config{
		GCTrigger: func() bool { return true },
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(c *Config)

func EnableTCO(enable bool) Option {
	return func(c *Config) {
		c.EnableTCO = enable
	}
}

func SetGCTrigger(trigger func() bool) Option {
	return func(c *Config) {
		c.GCTrigger = trigger
	}
}
