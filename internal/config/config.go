package config

import "os"

type ErrMissingEnv struct {
	envVar string
}

func (e *ErrMissingEnv) Error() string {
	return "missing env var: " + e.envVar
}

type ConfigOpt func(c *Config) error

// WithListenerAddr checks the env for LISTENER_ADDRESS
func WithListenerAddr() ConfigOpt {
	return func(c *Config) error {
		listenerAddr := os.Getenv("LISTENER_ADDRESS")
		if len(listenerAddr) == 0 {
			return &ErrMissingEnv{"LISTENER_ADDRESS"}
		}

		c.ListenerAddr = listenerAddr

		return nil
	}
}

// WithBrokers checks the env for BROKERS
func WithBrokers() ConfigOpt {
	return func(c *Config) error {
		brokers := os.Getenv("BROKERS")
		if len(brokers) == 0 {
			return &ErrMissingEnv{"BROKERS"}
		}

		c.Brokers = brokers

		return nil
	}
}

// WithOtelServiceName checks the env for OTEL_SERVICE_NAME
func WithOtelServiceName() ConfigOpt {
	return func(c *Config) error {
		serviceName := os.Getenv("OTEL_SERVICE_NAME")
		if len(serviceName) == 0 {
			return &ErrMissingEnv{"OTEL_SERVICE_NAME"}
		}

		c.OtelService = serviceName

		return nil
	}
}

// WithOtelEndpoint checks the env for OTEL_EXPORTER_OTLP_ENDPOINT
func WithOtelEndpoint() ConfigOpt {
	return func(c *Config) error {
		otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if len(otelEndpoint) == 0 {
			return &ErrMissingEnv{"OTEL_EXPORTER_OTLP_ENDPOINT"}
		}

		c.OtelEndpoint = otelEndpoint

		return nil
	}
}

// WithProfileURL checks the env for PROFILE_URL
func WithProfileURL() ConfigOpt {
	return func(c *Config) error {
		profileUrl := os.Getenv("PROFILE_URL")
		if len(profileUrl) == 0 {
			return &ErrMissingEnv{"PROFILE_URL"}
		}

		c.ProfileURL = profileUrl

		return nil
	}
}

// WithConsumerName checks the env for CONSUMER_GROUP_NAME
func WithConsumerName() ConfigOpt {
	return func(c *Config) error {
		name := os.Getenv("CONSUMER_GROUP_NAME")
		if len(name) == 0 {
			return &ErrMissingEnv{"CONSUMER_GROUP_NAME"}
		}

		c.ConsumerName = name

		return nil
	}
}

type Config struct {
	ListenerAddr string
	Brokers      string
	OtelService  string
	OtelEndpoint string
	ProfileURL   string
	ConsumerName string
}

func New(opts ...ConfigOpt) (*Config, error) {
	config := &Config{}

	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}
