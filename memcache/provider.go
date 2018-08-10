package memcache

import (
	"errors"
	"reflect"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/phachon/fasthttpsession"
)

// session MemCache provider

const ProviderName = "memcache"

var (
	provider = NewProvider()
	encrypt  = fasthttpsession.NewEncrypt()
)

type Provider struct {
	config         *Config
	values         *fasthttpsession.CCMap
	memCacheClient *memcache.Client
	maxLifeTime    int64
}

// new memcache provider
func NewProvider() *Provider {
	return &Provider{
		config:         &Config{},
		values:         fasthttpsession.NewDefaultCCMap(),
		memCacheClient: &memcache.Client{},
	}
}

// init provider config
func (mcp *Provider) Init(lifeTime int64, memCacheConfig fasthttpsession.ProviderConfig) error {
	if memCacheConfig.Name() != ProviderName {
		return errors.New("session memcache provider init error, config must memcache config")
	}
	vc := reflect.ValueOf(memCacheConfig)
	rc := vc.Interface().(*Config)
	mcp.config = rc

	// config check
	if len(mcp.config.ServerList) == 0 {
		return errors.New("session memcache provider init error, config ServerList not empty")
	}
	if mcp.config.MaxIdle <= 0 {
		return errors.New("session memcache provider init error, config MaxIdle must be more than 0")
	}
	// init config serialize func
	if mcp.config.SerializeFunc == nil {
		mcp.config.SerializeFunc = encrypt.GobEncode
	}
	if mcp.config.UnSerializeFunc == nil {
		mcp.config.UnSerializeFunc = encrypt.GobDecode
	}
	// create memcache client
	mcp.memCacheClient = memcache.New(mcp.config.ServerList...)
	mcp.memCacheClient.MaxIdleConns = mcp.config.MaxIdle
	mcp.maxLifeTime = lifeTime
	return nil
}

// not need gc
func (mcp *Provider) NeedGC() bool {
	return false
}

// session memcache provider not need garbage collection
func (mcp *Provider) GC() {}

// read session store by session id
func (mcp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	memClient := mcp.getMemCacheClient()

	item, err := memClient.Get(mcp.getMemCacheSessionKey(sessionId))
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return NewMemCacheStore(sessionId), nil
		} else {
			return nil, err
		}
	}
	if len(item.Value) == 0 {
		return NewMemCacheStore(sessionId), nil
	}

	data, err := mcp.config.UnSerializeFunc(item.Value)
	if err != nil {
		return nil, err
	}

	return NewMemCacheStoreData(sessionId, data), nil
}

// regenerate session
func (mcp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	memClient := mcp.getMemCacheClient()

	item, err := memClient.Get(mcp.getMemCacheSessionKey(oldSessionId))
	if err != nil || len(item.Value) == 0 {
		// false, old sessionId not exists
		err := memClient.Set(&memcache.Item{
			Key:        mcp.getMemCacheSessionKey(sessionId),
			Value:      []byte(""),
			Expiration: int32(mcp.maxLifeTime),
		})
		if err != nil {
			return nil, err
		}
		return NewMemCacheStore(sessionId), nil
	}
	// true, old sessionId exists, delete old sessionId
	err = memClient.Delete(mcp.getMemCacheSessionKey(oldSessionId))
	if err != nil {
		return nil, err
	}
	item.Key = mcp.getMemCacheSessionKey(sessionId)
	item.Expiration = int32(mcp.maxLifeTime)
	err = memClient.Set(item)
	if err != nil {
		return nil, err
	}

	return mcp.ReadStore(sessionId)
}

// destroy session by sessionId
func (mcp *Provider) Destroy(sessionId string) error {
	memClient := mcp.getMemCacheClient()
	return memClient.Delete(mcp.getMemCacheSessionKey(sessionId))
}

// session values count
func (mcp *Provider) Count() int {
	return 0
}

// get memcache session key, prefix:sessionId
func (mcp *Provider) getMemCacheSessionKey(sessionId string) string {
	return mcp.config.KeyPrefix + ":" + sessionId
}

func (mcp *Provider) getMemCacheClient() *memcache.Client {
	if mcp.memCacheClient == nil {
		mcp.memCacheClient = memcache.New(mcp.config.ServerList...)
	}
	return mcp.memCacheClient
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, provider)
}
