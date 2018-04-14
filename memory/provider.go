package memory

import (
	"github.com/phachon/fasthttpsession"
	"time"
	"errors"
	"reflect"
)

// session memory provider

const ProviderName = "memory"

type Provider struct {
	config *Config
	values *fasthttpsession.CCMap
}

// new memory provider
func NewProvider() *Provider {
	return &Provider{
		config: &Config{},
		values: fasthttpsession.NewDefaultCCMap(),
	}
}

// init provider config
func (mp *Provider) Init(memoryConfig fasthttpsession.ProviderConfig) error {
	if memoryConfig.Name() != ProviderName {
		return errors.New("session memory provider init error, config must memory config")
	}
	vc := reflect.ValueOf(memoryConfig)
	mc := vc.Interface().(*Config)
	mp.config = mc
	return nil
}

// session garbage collection
func (mp *Provider) GC(sessionLifetime int64) {
	for sessionId, value := range mp.values.GetAll() {
		if time.Now().Unix() >= value.(*Store).lastActiveTime + sessionLifetime {
			// destroy session sessionId
			mp.Destroy(sessionId)
			return
		}
	}
}

// session id is exist
func (mp *Provider) SessionIdIsExist(sessionId string) bool {
	ok := mp.values.IsExist(sessionId)
	if ok {
		return true
	}
	return false
}

// read session store by session id
func (mp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {
	memStore := mp.values.Get(sessionId)
	if memStore != nil {
		return memStore.(*Store), nil
	}

	newMemStore := NewMemoryStore(sessionId)
	mp.values.Set(sessionId, newMemStore)

	return newMemStore, nil
}

// regenerate session
func (mp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {
	memStoreInter := mp.values.Get(oldSessionId)
	if memStoreInter != nil {
		memStore := memStoreInter.(*Store)
		// insert new session store
		newMemStore := NewMemoryStoreData(sessionId, memStore.GetAll())
		mp.values.Set(sessionId, newMemStore)
		// delete old session store
		mp.values.Delete(oldSessionId)
		return newMemStore, nil
	}

	memStore := NewMemoryStore(sessionId)
	mp.values.Set(sessionId, memStore)

	return memStore, nil
}

// destroy session by sessionId
func (mp *Provider) Destroy(sessionId string) error {
	mp.values.Delete(sessionId)
	return nil
}

// session values count
func (mp *Provider) Count() int {
	return mp.values.Count()
}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, NewProvider())
}