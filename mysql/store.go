package mysql

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
	"time"
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

	b, err := utils.GobEncode(ms.GetAll())
	if err != nil {
		return err
	}
	_, err = provider.sessionDao.updateBySessionId(ms.GetSessionId(), string(b), time.Now().Unix())
	return err
}