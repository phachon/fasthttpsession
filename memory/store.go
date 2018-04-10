package memory

import (
	"fasthttpsession"
	"github.com/valyala/fasthttp"
	"time"
)

// session memory store

func NewMemoryStore() *Store {
	memStore := &Store{}
	memStore.SessionId = ""
	memStore.Data = make(map[interface{}]interface{})
	memStore.lastActiveTime = time.Now().Unix()
	return memStore
}

type Store struct {
	fasthttpsession.Store
	lastActiveTime int64
}

// save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	return nil
}
