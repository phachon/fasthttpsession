package sqlite3

import (
	"time"

	"github.com/savsgio/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session sqlite3 store

// new default sqlite3 store
func NewSqLite3Store(sessionId string) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionId, make(map[string]interface{}))
	return sqlite3Store
}

// new sqlite3 store data
func NewSqLite3StoreData(sessionId string, data map[string]interface{}) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionId, data)
	return sqlite3Store
}

type Store struct {
	fasthttpsession.Store
}

// save store
func (ss *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ss.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionId(ss.GetSessionId())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionId(ss.GetSessionId(), string(b), time.Now().Unix())
	return err
}
