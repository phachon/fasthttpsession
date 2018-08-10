package mysql

import (
	"time"

	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session mysql store

// new default mysql store
func NewMysqlStore(sessionId string) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionId, make(map[string]interface{}))
	return mysqlStore
}

// new mysql store data
func NewMysqlStoreData(sessionId string, data map[string]interface{}) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionId, data)
	return mysqlStore
}

type Store struct {
	fasthttpsession.Store
}

// save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := provider.config.SerializeFunc(ms.GetAll())
	if err != nil {
		return err
	}
	session, err := provider.sessionDao.getSessionBySessionId(ms.GetSessionId())
	if err != nil || len(session) == 0 {
		return nil
	}
	_, err = provider.sessionDao.updateBySessionId(ms.GetSessionId(), string(b), time.Now().Unix())
	return err
}
