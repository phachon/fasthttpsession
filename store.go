package fasthttpsession

import (
	"sync"
	"github.com/valyala/fasthttp"
	"time"
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
	SessionId       string
	Lock            sync.RWMutex
	Data            map[interface{}]interface{}
	LastActiveTime  int64
}

// get Data by key
func (s *Store) Get(key interface{}) interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	value, ok := s.Data[key]
	if ok {
		return value
	}
	return nil
}

// get all Data
func (s *Store) GetAll() map[interface{}]interface{} {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	return s.Data
}

// set Data
func (s *Store) Set(key interface{}, value interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data[key] = value
}

// delete Data by key
func (s *Store) Delete(key interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Data, key)
}

// flush all Data
func (s *Store) Flush() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data = make(map[interface{}]interface{})
}

// get session id
func (s *Store) GetSessionId() string {
	return s.SessionId
}

// update session lastActiveTime
func (s *Store) UpdateLastActiveTime() {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.LastActiveTime = time.Now().Unix()
}
