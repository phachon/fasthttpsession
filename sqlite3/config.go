package sqlite3

// session sqlite3 config

type Config struct {

	// sqlite3 db file path
	DBPath string

	// session table name
	TableName string

	// sqlite3 max free idle
	SetMaxIdleConn int

	// sqlite3 max open idle
	SetMaxOpenConn int

	// session value serialize func
	SerializeFunc func(data map[string]interface{}) ([]byte, error)

	// session value unSerialize func
	UnSerializeFunc func(data []byte) (map[string]interface{}, error)
}

func NewConfigWith(dbPath, tableName string) (cf *Config) {
	cf = &Config{
		SetMaxOpenConn: 500,
		SetMaxIdleConn: 50,
	}
	cf.DBPath = dbPath
	cf.TableName = tableName
	return
}

func (sc *Config) Name() string {
	return ProviderName
}