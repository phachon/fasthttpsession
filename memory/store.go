package memory

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
	"time"
	"sync"
)

// session memory store

// new default memory store
func NewMemoryStore(sessionId string) *Store {
	memStore := &Store{}
	memStore.Init(sessionId, make(map[string]interface{}))
	return memStore
}

// new memory store data
func NewMemoryStoreData(sessionId string, data map[string]interface{}) *Store {
	memStore := &Store{}
	memStore.Init(sessionId, data)
	return memStore
}

type Store struct {
	fasthttpsession.Store
	lock            sync.RWMutex
	lastActiveTime  int64
}

// save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.lastActiveTime = time.Now().Unix()
	return nil
}