package config

import (
	"os"
	"strconv"
)

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

// WithUserEventsTopic checks the env for USER_EVENTS_TOPIC
func WithUserEventsTopic() ConfigOpt {
	return func(c *Config) error {
		topic := os.Getenv("USER_EVENTS_TOPIC")
		if len(topic) == 0 {
			return &ErrMissingEnv{"USER_EVENTS_TOPIC"}
		}

		c.UserEventsTopic = topic

		return nil
	}
}

// WithWSAddress checks the env for WS_ADDRESS
func WithWSAddress() ConfigOpt {
	return func(c *Config) error {
		address := os.Getenv("WS_ADDRESS")
		if len(address) == 0 {
			return &ErrMissingEnv{"WS_ADDRESS"}
		}

		c.WSAddress = address

		return nil
	}
}

// WithMaxWSClients checks the env for MAX_WS_CLIENTS
func WithMaxWSClients() ConfigOpt {
	return func(c *Config) error {
		maxClients := os.Getenv("MAX_WS_CLIENTS")
		if len(maxClients) == 0 {
			return &ErrMissingEnv{"MAX_WS_CLIENTS"}
		}

		maxClientsInt, err := strconv.Atoi(maxClients)
		if err != nil {
			return err
		}

		c.MaxWSClients = maxClientsInt

		return nil
	}
}

type Config struct {
	ListenerAddr    string
	Brokers         string
	OtelService     string
	OtelEndpoint    string
	ProfileURL      string
	ConsumerName    string
	UserEventsTopic string
	WSAddress       string
	MaxWSClients    int
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
