package file

// session file config

type Config struct {
	SavePath string
	Suffix   string
	SerializeFunc func(data map[string]interface{}) ([]byte, error)
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

func (fc *Config) Name() string {
	return ProviderName
}