package sqlite3

import (
	"errors"
	"github.com/phachon/fasthttpsession"
	"reflect"
	"time"
)

// session sqlite3 provider

//  session Table structure
//
//  DROP TABLE IF EXISTS `session`;
//  CREATE TABLE `session` (
//    `session_id` varchar(64) NOT NULL DEFAULT '',
//    `contents` TEXT NOT NULL,
//    `last_active` int(10) NOT NULL DEFAULT '0',
//    PRIMARY KEY (`session_id`),
//  )
//  create index last_active on session (last_active);
//

const ProviderName = "sqlite3"

var (
	provider = NewProvider()
	encrypt  = fasthttpsession.NewEncrypt()
)

type Provider struct {
	config      *Config
	values      *fasthttpsession.CCMap
	sessionDao  *sessionDao
	maxLifeTime int64
}

// new sqlite3 provider
func NewProvider() *Provider {
	return &Provider{
		config:     &Config{},
		values:     fasthttpsession.NewDefaultCCMap(),
		sessionDao: &sessionDao{},
	}
}

// init provider config
func (sp *Provider) Init(lifeTime int64, sqlite3Config fasthttpsession.ProviderConfig) error {
	if sqlite3Config.Name() != ProviderName {
		return errors.New("session sqlite3 provider init error, config must sqlite3 config")
	}
	vc := reflect.ValueOf(sqlite3Config)
	rc := vc.Interface().(*Config)
	sp.config = rc
	sp.maxLifeTime = lifeTime

	// check config
	if sp.config.DBPath == "" {
		return errors.New("session sqlite3 provider init error, config DBPath not empty")
	}
	// init config serialize func
	if sp.config.SerializeFunc == nil {
		sp.config.SerializeFunc = encrypt.Base64Encode
	}
	if sp.config.UnSerializeFunc == nil {
		sp.config.UnSerializeFunc = encrypt.Base64Decode
	}
	// init sessionDao
	sessionDao, err := newSessionDao(sp.config.DBPath, sp.config.TableName)
	if err != nil {
		return err
	}
	sessionDao.sqlite3Conn.SetMaxOpenConns(sp.config.SetMaxIdleConn)
	sessionDao.sqlite3Conn.SetMaxIdleConns(sp.config.SetMaxIdleConn)

	sp.sessionDao = sessionDao
	return sessionDao.sqlite3Conn.Ping()
}

// not need gc
func (sp *Provider) NeedGC() bool {
	return true
}

// session sqlite3 provider not need garbage collection
func (sp *Provider) GC() {
	sp.sessionDao.deleteSessionByMaxLifeTime(sp.maxLifeTime)
}

// read session store by session id
func (sp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	sessionValue, err := sp.sessionDao.getSessionBySessionId(sessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		_, err := sp.sessionDao.insert(sessionId, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewSqLite3Store(sessionId), nil
	}
	if len(sessionValue["contents"]) == 0 {
		return NewSqLite3Store(sessionId), nil
	}

	data, err := sp.config.UnSerializeFunc(sessionValue["contents"])
	if err != nil {
		return nil, err
	}

	return NewSqLite3StoreData(sessionId, data), nil
}

// regenerate session
func (sp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	sessionValue, err := sp.sessionDao.getSessionBySessionId(oldSessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		// old sessionId not exists, insert new sessionId
		_, err := sp.sessionDao.insert(sessionId, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewSqLite3Store(sessionId), nil
	}

	// delete old session
	_, err = sp.sessionDao.deleteBySessionId(oldSessionId)
	if err != nil {
		return nil, err
	}
	// insert new session
	_, err = sp.sessionDao.insert(sessionId, string(sessionValue["contents"]), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return sp.ReadStore(sessionId)
}

// destroy session by sessionId
func (sp *Provider) Destroy(sessionId string) error {
	_, err := sp.sessionDao.deleteBySessionId(sessionId)
	return err
}

// session values count
func (sp *Provider) Count() int {
	return sp.sessionDao.countSessions()
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, provider)
}
