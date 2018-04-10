package memory

import (
	"fasthttpsession"
	"sync"
)

// session memory provider

const ProviderName = "memory"

type Provider struct {
	lock sync.RWMutex
	values map[string]*Store
}

func NewProvider() *Provider {
	return &Provider{
		values: make(map[string]*Store),
	}
}

// init provider config
func (mp *Provider) Init(memoryConfig fasthttpsession.ProviderConfig) error {
	return nil
}

// session id is exist
func (mp *Provider) SessionIdIsExist(sessionId string) bool {
	mp.lock.RLock()
	defer mp.lock.RUnlock()
	_, ok := mp.values[sessionId]
	if ok {
		return true
	}
	return false
}

// read session store by session id
func (mp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {
	mp.lock.RLock()
	memStore, ok := mp.values[sessionId]
	if ok {
		mp.lock.RUnlock()
		return memStore, nil
	}
	mp.lock.RUnlock()

	mp.lock.Lock()
	mp.values[sessionId] = NewMemoryStore()
	mp.lock.Unlock()

	return NewMemoryStore(), nil
}

// destroy session by sessionId
func (mp *Provider) Destroy(sessionId string) error {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	delete(mp.values, sessionId)
	return nil
}

// session values count
func (mp *Provider) Count() int {
	return len(mp.values)
}

func init()  {
	fasthttpsession.Register(ProviderName, NewProvider())
}