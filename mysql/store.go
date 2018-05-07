package mysql

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session mysql store

// new default mysql store
func NewmysqlStore(sessionId string) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionId, make(map[string]interface{}))
	return mysqlStore
}

// new mysql store data
func NewmysqlStoreData(sessionId string, data map[string]interface{}) *Store {
	mysqlStore := &Store{}
	mysqlStore.Init(sessionId, data)
	return mysqlStore
}

type Store struct {
	fasthttpsession.Store
}

// save store
func (rs *Store) Save(ctx *fasthttp.RequestCtx) error {

	b, err := utils.GobEncode(rs.GetAll())
	if err != nil {
		return err
	}
	conn := provider.mysqlPool.Get()
	defer conn.Close()
	conn.Do("SETEX", provider.getmysqlSessionKey(rs.GetSessionId()), provider.maxLifeTime, string(b))

	return nil
}