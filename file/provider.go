package file

import (
	"github.com/phachon/fasthttpsession"
	"sync"
	"errors"
	"reflect"
	"os"
	"path"
)

// session file provider

const ProviderName = "file"

type Provider struct {
	lock sync.Mutex
	file *file
	config *Config
}

// new file provider
func NewProvider() *Provider {
	return &Provider{
		file: &file{},
		config: &Config{},
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

	// path is exist
	if !fp.file.pathIsExists(fp.config.SavePath) {
		return os.MkdirAll(fp.config.SavePath, 0777)
	}

	return nil
}

// session garbage collection
func (fp *Provider) GC(sessionLifetime int64) {

}

// session id is exist
func (fp *Provider) SessionIdIsExist(sessionId string) bool {
	fp.lock.Lock()
	defer fp.lock.Unlock()

	return fp.file.pathIsExists(fp.getFilePath(sessionId))
}

// read session store by session id
func (fp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {
	fp.lock.Lock()
	defer fp.lock.Unlock()

	//sessionInfo, err := fp.file.getContent(fp.getFilePath(sessionId))
	//if err != nil {
	//
	//} else {
	//
	//}
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

// get session file path
func (fp *Provider) getFilePath(sessionId string) string {
	return path.Join(fp.config.SavePath, string(sessionId[0]), string(sessionId[1]), fp.sessionFileName(sessionId))
}

// get sessionId filename
func (fp *Provider) sessionFileName(sessionId string) string {
	return sessionId + fp.config.Suffix
}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, NewProvider())
}