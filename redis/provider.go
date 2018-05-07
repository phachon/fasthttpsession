package redis

import (
	"github.com/phachon/fasthttpsession"
	"errors"
	"reflect"
	"github.com/gomodule/redigo/redis"
)

// session redis provider

const ProviderName = "redis"

var (
	utils = fasthttpsession.NewUtils()
	provider = NewProvider()
)

type Provider struct {
	config *Config
	values *fasthttpsession.CCMap
	redisPool *redis.Pool
	maxLifeTime int64
}

// new redis provider
func NewProvider() *Provider {
	return &Provider{
		config: &Config{},
		values: fasthttpsession.NewDefaultCCMap(),
		redisPool: &redis.Pool{},
	}
}

// init provider config
func (rp *Provider) Init(redisConfig fasthttpsession.ProviderConfig) error {
	if redisConfig.Name() != ProviderName {
		return errors.New("session redis provider init error, config must redis config")
	}
	vc := reflect.ValueOf(redisConfig)
	rc := vc.Interface().(*Config)
	rp.config = rc

	// create redis conn pool
	rp.redisPool = newRedisPool(rp.config)

	// check redis conn
	conn := rp.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return err
	}
	return nil
}

// set maxLifeTime
func (rp *Provider) MaxLifeTime(lifeTime int64)  {
	rp.maxLifeTime = lifeTime
}

// not need gc
func (rp *Provider) NeedGC() bool {
	return false
}

// session redis provider not need garbage collection
func (rp *Provider) GC(sessionLifetime int64) {}


// read session store by session id
func (rp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	conn := rp.redisPool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", rp.getRedisSessionKey(sessionId)))
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if len(reply) == 0 {
		conn.Do("SET", rp.getRedisSessionKey(sessionId), "", "EX", rp.maxLifeTime)
		return NewRedisStore(sessionId), nil
	}

	data, err := utils.GobDecode(reply)
	if err != nil {
		return nil, err
	}

	return NewRedisStoreData(sessionId, data), nil
}

// regenerate session
func (rp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	conn := rp.redisPool.Get()
	defer conn.Close()

	existed, err := redis.Int(conn.Do("EXISTS", rp.getRedisSessionKey(oldSessionId)))
	if err != nil || existed == 0 {
		// false
		conn.Do("SET", rp.getRedisSessionKey(sessionId), "", "EX", rp.maxLifeTime)
		return NewRedisStore(sessionId), nil
	}
	// true
	conn.Do("RENAME", rp.getRedisSessionKey(oldSessionId), rp.getRedisSessionKey(sessionId))
	conn.Do("EXPIRE", rp.getRedisSessionKey(sessionId), rp.maxLifeTime)

	return rp.ReadStore(sessionId)
}

// destroy session by sessionId
func (rp *Provider) Destroy(sessionId string) error {
	conn := rp.redisPool.Get()
	defer conn.Close()

	existed, err := redis.Int(conn.Do("EXISTS", rp.getRedisSessionKey(sessionId)))
	if err != nil || existed == 0 {
		return nil
	}
	conn.Do("DEL", rp.getRedisSessionKey(sessionId))
	return nil
}

// session values count
func (rp *Provider) Count() int {
	conn := rp.redisPool.Get()
	defer conn.Close()

	replyMap, err := redis.Strings(conn.Do("KEYS", rp.config.KeyPrefix+":*"))
	if err != nil {
		return 0
	}
	return len(replyMap)
}

// get redis session key, prefix:sessionId
func (rp *Provider) getRedisSessionKey(sessionId string) string {
	return rp.config.KeyPrefix+":"+sessionId
}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, provider)
}