package file

// session file config

type Config struct {

	// session file save path
	SavePath string

	// session file suffix
	Suffix string

	// session value serialize func
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

func (fc *Config) Name() string {
	return ProviderName
}
