package redis

import (
	"errors"
	"reflect"

	"github.com/gomodule/redigo/redis"
	"github.com/phachon/fasthttpsession"
)

// session redis provider

const ProviderName = "redis"

var (
	provider = NewProvider()
	encrypt  = fasthttpsession.NewEncrypt()
)

type Provider struct {
	config      *Config
	values      *fasthttpsession.CCMap
	redisPool   *redis.Pool
	maxLifeTime int64
}

// new redis provider
func NewProvider() *Provider {
	return &Provider{
		config:    &Config{},
		values:    fasthttpsession.NewDefaultCCMap(),
		redisPool: &redis.Pool{},
	}
}

// init provider config
func (rp *Provider) Init(lifeTime int64, redisConfig fasthttpsession.ProviderConfig) error {
	if redisConfig.Name() != ProviderName {
		return errors.New("session redis provider init error, config must redis config")
	}
	vc := reflect.ValueOf(redisConfig)
	rc := vc.Interface().(*Config)
	rp.config = rc
	rp.maxLifeTime = lifeTime

	// config check
	if rp.config.Host == "" {
		return errors.New("session redis provider init error, config Host not empty")
	}
	if rp.config.Port == 0 {
		return errors.New("session redis provider init error, config Port not empty")
	}
	if rp.config.MaxIdle <= 0 {
		return errors.New("session redis provider init error, config MaxIdle must be more than 0")
	}
	if rp.config.IdleTimeout <= 0 {
		return errors.New("session redis provider init error, config IdleTimeout must be more than 0")
	}
	// init config serialize func
	if rp.config.SerializeFunc == nil {
		rp.config.SerializeFunc = encrypt.GobEncode
	}
	if rp.config.UnSerializeFunc == nil {
		rp.config.UnSerializeFunc = encrypt.GobDecode
	}
	// create redis conn pool
	rp.redisPool = newRedisPool(rp.config)

	// check redis conn
	conn := rp.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return errors.New("session redis provider init error, " + err.Error())
	}
	return nil
}

// not need gc
func (rp *Provider) NeedGC() bool {
	return false
}

// session redis provider not need garbage collection
func (rp *Provider) GC() {}

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

	data, err := rp.config.UnSerializeFunc(reply)
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
	return rp.config.KeyPrefix + ":" + sessionId
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, provider)
}
