package redis

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// session redis store

// new default redis store
func NewRedisStore(sessionId string) *Store {
	redisStore := &Store{}
	redisStore.Init(sessionId, make(map[string]interface{}))
	return redisStore
}

// new redis store data
func NewRedisStoreData(sessionId string, data map[string]interface{}) *Store {
	redisStore := &Store{}
	redisStore.Init(sessionId, data)
	return redisStore
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
	conn := provider.redisPool.Get()
	defer conn.Close()
	conn.Do("SETEX", provider.getRedisSessionKey(rs.GetSessionId()), provider.maxLifeTime, string(b))

	return nil
}