package file

// session file config

type Config struct {
	SavePath string
}

func (fc *Config) Name() string {
	return ProviderName
}