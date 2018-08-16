package memory

import (
	"errors"
	"reflect"
	"time"

	"github.com/savsgio/fasthttpsession"
)

// session memory provider

const ProviderName = "memory"

type Provider struct {
	config      *Config
	values      *fasthttpsession.CCMap
	maxLifeTime int64
}

// new memory provider
func NewProvider() *Provider {
	return &Provider{
		config:      &Config{},
		values:      fasthttpsession.NewDefaultCCMap(),
		maxLifeTime: 0,
	}
}

// init provider config
func (mp *Provider) Init(lifeTime int64, memoryConfig fasthttpsession.ProviderConfig) error {
	if memoryConfig.Name() != ProviderName {
		return errors.New("session memory provider init error, config must memory config")
	}
	vc := reflect.ValueOf(memoryConfig)
	mc := vc.Interface().(*Config)
	mp.config = mc

	mp.maxLifeTime = lifeTime
	return nil
}

// need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// session garbage collection
func (mp *Provider) GC() {
	for sessionID, value := range mp.values.GetAll() {
		if time.Now().Unix() >= value.(*Store).lastActiveTime+mp.maxLifeTime {
			// destroy session sessionID
			mp.Destroy(sessionID)
			return
		}
	}
}

// read session store by session id
func (mp *Provider) ReadStore(sessionID string) (fasthttpsession.SessionStore, error) {
	memStore := mp.values.Get(sessionID)
	if memStore != nil {
		return memStore.(*Store), nil
	}

	newMemStore := NewMemoryStore(sessionID)
	mp.values.Set(sessionID, newMemStore)

	return newMemStore, nil
}

// regenerate session
func (mp *Provider) Regenerate(oldSessionId string, sessionID string) (fasthttpsession.SessionStore, error) {
	memStoreInter := mp.values.Get(oldSessionId)
	if memStoreInter != nil {
		memStore := memStoreInter.(*Store)
		// insert new session store
		newMemStore := NewMemoryStoreData(sessionID, memStore.GetAll())
		mp.values.Set(sessionID, newMemStore)
		// delete old session store
		mp.values.Delete(oldSessionId)
		return newMemStore, nil
	}

	memStore := NewMemoryStore(sessionID)
	mp.values.Set(sessionID, memStore)

	return memStore, nil
}

// destroy session by sessionID
func (mp *Provider) Destroy(sessionID string) error {
	mp.values.Delete(sessionID)
	return nil
}

// session values count
func (mp *Provider) Count() int {
	return mp.values.Count()
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, NewProvider())
}
