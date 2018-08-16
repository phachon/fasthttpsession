package memory

import (
	"sync"
	"time"

	"github.com/savsgio/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session memory store

// NewMemoryStore new default memory store
func NewMemoryStore(sessionID string) *Store {
	memStore := &Store{}
	memStore.Init(sessionID, make(map[string]interface{}))
	return memStore
}

// NewMemoryStoreData new memory store data
func NewMemoryStoreData(sessionID string, data map[string]interface{}) *Store {
	memStore := &Store{}
	memStore.Init(sessionID, data)
	return memStore
}

// Store store struct
type Store struct {
	fasthttpsession.Store
	lock           sync.RWMutex
	lastActiveTime int64
}

// Save save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.lastActiveTime = time.Now().Unix()
	return nil
}
