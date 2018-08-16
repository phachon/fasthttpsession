package postgres

import (
	"time"

	"github.com/savsgio/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session postgres store

// new default postgres store
func NewPostgresStore(sessionId string) *Store {
	postgresStore := &Store{}
	postgresStore.Init(sessionId, make(map[string]interface{}))
	return postgresStore
}

// NewPostgresStoreData new postgres store data
func NewPostgresStoreData(sessionId string, data map[string]interface{}) *Store {
	postgresStore := &Store{}
	postgresStore.Init(sessionId, data)
	return postgresStore
}

type Store struct {
	fasthttpsession.Store
}

// Save save store
func (ps *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ps.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionId(ps.GetSessionId())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionId(ps.GetSessionId(), string(b), time.Now().Unix())
	return err
}
