package memory

// session memory config

type Config struct {
}

func (mc *Config) Name() string {
	return ProviderName
}
