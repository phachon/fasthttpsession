package sqlite3

import (
	"time"

	"github.com/savsgio/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session sqlite3 store

// NewSqLite3Store new default sqlite3 store
func NewSqLite3Store(sessionID string) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionID, make(map[string]interface{}))
	return sqlite3Store
}

// NewSqLite3StoreData new sqlite3 store data
func NewSqLite3StoreData(sessionID string, data map[string]interface{}) *Store {
	sqlite3Store := &Store{}
	sqlite3Store.Init(sessionID, data)
	return sqlite3Store
}

// Store store struct
type Store struct {
	fasthttpsession.Store
}

// Save save store
func (ss *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ss.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionID(ss.GetSessionID())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionID(ss.GetSessionID(), string(b), time.Now().Unix())
	return err
}
