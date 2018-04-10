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
	SessionId    string
	lock   sync.RWMutex
	Data map[interface{}]interface{}
}

// get Data by key
func (s *Store) Get(key interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.Data[key]
	if ok {
		return value
	}
	return nil
}

// get all Data
func (s *Store) GetAll() map[interface{}]interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Data
}

// set Data
func (s *Store) Set(key interface{}, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data[key] = value
}

// delete Data by key
func (s *Store) Delete(key interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.Data, key)
}

// flush all Data
func (s *Store) Flush() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data = make(map[interface{}]interface{})
}

// get session id
func (s *Store) GetSessionId() string {
	return s.SessionId
}
