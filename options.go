package censor

// OptsConfig is a options configuration for the Processor.
type OptsConfig struct {
	config     *Config
	configPath string
}

// Option is a function that sets some option on the Processor.
type Option func(*OptsConfig)

// WithConfig returns an Option that sets the configuration on the Processor.
func WithConfig(c *Config) func(*OptsConfig) {
	return func(o *OptsConfig) {
		o.config = c
	}
}

// WithConfigPath returns an Option that sets the configuration path on the Processor.
func WithConfigPath(path string) func(*OptsConfig) {
	return func(o *OptsConfig) {
		o.configPath = path
	}
}
