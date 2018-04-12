package file

import (
	"fasthttpsession"
	"sync"
	"errors"
	"reflect"
)

// session file provider

const ProviderName = "file"

type Provider struct {
	lock sync.Mutex
	config *Config
}

// new file provider
func NewProvider() *Provider {
	return &Provider{
		config: &Config{
			SavePath: "",
		},
	}
}

// init provider config
func (fp *Provider) Init(fileConfig fasthttpsession.ProviderConfig) error {
	if fileConfig.Name() != ProviderName {
		return errors.New("session file provider init error, config must file config")
	}

	vc := reflect.ValueOf(fileConfig)
	fc := vc.Interface().(*Config)
	fp.config = fc
	return nil
}

// session garbage collection
func (fp *Provider) GC(sessionLifetime int64) {

}

// session id is exist
func (fp *Provider) SessionIdIsExist(sessionId string) bool {
	fp.lock.Lock()
	defer fp.lock.Unlock()


	return false
}

// read session store by session id
func (fp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {


	return &Store{}, nil
}



// regenerate session
func (fp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	return &Store{}, nil
}

// destroy session by sessionId
func (fp *Provider) Destroy(sessionId string) error {
	return nil
}

// session values count
func (fp *Provider) Count() int {
	return 0
}

// get sessionId file
func (fp *Provider) GetSessionFile(sessionId string) {

}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, NewProvider())
}