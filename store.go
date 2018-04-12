package fasthttpsession

import (
	"sync"
	"github.com/valyala/fasthttp"
)

// session store struct

type SessionStore interface {
	Save(*fasthttp.RequestCtx) error
	Get(key interface{}) interface{}
	GetAll() map[interface{}]interface{}
	Set(key interface{}, value interface{})
	Delete(key interface{})
	Flush()
	GetSessionId() string
}

type Store struct {
	sessionId       string
	Lock            sync.RWMutex
	data            map[interface{}]interface{}
}

// init store data and sessionId
func (s *Store) Init(sessionId string, data map[interface{}]interface{}) {
	s.sessionId = sessionId
	s.data = data
}

// get data by key
func (s *Store) Get(key interface{}) interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	value, ok := s.data[key]
	if ok {
		return value
	}
	return nil
}

// get all data
func (s *Store) GetAll() map[interface{}]interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.data
}

// set data
func (s *Store) Set(key interface{}, value interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.data[key] = value
}

// delete data by key
func (s *Store) Delete(key interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.data, key)
}

// flush all data
func (s *Store) Flush() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.data = make(map[interface{}]interface{})
}

// get session id
func (s *Store) GetSessionId() string {
	return s.sessionId
}

// lock
func (s *Store) LockHandle(handle func()) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	handle()
}