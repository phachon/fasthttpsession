package memory

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
	"time"
)

// session memory store

func NewMemoryStore() *Store {
	memStore := &Store{}
	memStore.SessionId = ""
	memStore.Data = make(map[interface{}]interface{})
	memStore.LastActiveTime = time.Now().Unix()
	return memStore
}

type Store struct {
	fasthttpsession.Store
}

// save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.UpdateLastActiveTime()
	return nil
}