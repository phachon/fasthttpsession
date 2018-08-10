package fasthttpsession

import (
	"github.com/valyala/fasthttp"
)

// session store struct

type SessionStore interface {
	Save(*fasthttp.RequestCtx) error
	Get(key string) interface{}
	GetAll() map[string]interface{}
	Set(key string, value interface{})
	Delete(key string)
	Flush()
	GetSessionId() string
}

type Store struct {
	sessionId string
	data      *CCMap
}

// init store data and sessionId
func (s *Store) Init(sessionId string, data map[string]interface{}) {
	s.sessionId = sessionId
	s.data = NewDefaultCCMap()
	s.data.MSet(data)
}

// get data by key
func (s *Store) Get(key string) interface{} {
	return s.data.Get(key)
}

// get all data
func (s *Store) GetAll() map[string]interface{} {
	return s.data.GetAll()
}

// set data
func (s *Store) Set(key string, value interface{}) {
	s.data.Set(key, value)
}

// delete data by key
func (s *Store) Delete(key string) {
	s.data.Delete(key)
}

// flush all data
func (s *Store) Flush() {
	s.data.Clear()
}

// get session id
func (s *Store) GetSessionId() string {
	return s.sessionId
}
