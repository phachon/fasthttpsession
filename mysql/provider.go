package mysql

import (
	"errors"
	"github.com/phachon/fasthttpsession"
	"reflect"
	"time"
)

// session mysql provider

// session Table structure
//
// DROP TABLE IF EXISTS `session`;
// CREATE TABLE `session` (
//    `session_id` varchar(64) NOT NULL DEFAULT '' COMMENT 'Session id',
//    `contents` TEXT NOT NULL COMMENT 'Session data',
//    `last_active` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Last active time',
//    PRIMARY KEY (`session_id`),
//    KEY `last_active` (`last_active`)
// ) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='session table';
//

const ProviderName = "mysql"

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

// new mysql provider
func NewProvider() *Provider {
	return &Provider{
		config:     &Config{},
		values:     fasthttpsession.NewDefaultCCMap(),
		sessionDao: &sessionDao{},
	}
}

// init provider config
func (mp *Provider) Init(lifeTime int64, mysqlConfig fasthttpsession.ProviderConfig) error {
	if mysqlConfig.Name() != ProviderName {
		return errors.New("session mysql provider init error, config must mysql config")
	}
	vc := reflect.ValueOf(mysqlConfig)
	rc := vc.Interface().(*Config)
	mp.config = rc
	mp.maxLifeTime = lifeTime

	// check config
	if mp.config.Host == "" {
		return errors.New("session mysql provider init error, config Host not empty")
	}
	if mp.config.Port == 0 {
		return errors.New("session mysql provider init error, config Port not empty")
	}
	// init config serialize func
	if mp.config.SerializeFunc == nil {
		mp.config.SerializeFunc = encrypt.Base64Encode
	}
	if mp.config.UnSerializeFunc == nil {
		mp.config.UnSerializeFunc = encrypt.Base64Decode
	}
	// init sessionDao
	sessionDao, err := newSessionDao(mp.config.getMysqlDSN(), mp.config.TableName)
	if err != nil {
		return err
	}
	sessionDao.mysqlConn.SetMaxOpenConns(mp.config.SetMaxIdleConn)
	sessionDao.mysqlConn.SetMaxIdleConns(mp.config.SetMaxIdleConn)

	mp.sessionDao = sessionDao
	return sessionDao.mysqlConn.Ping()
}

// not need gc
func (mp *Provider) NeedGC() bool {
	return true
}

// session mysql provider not need garbage collection
func (mp *Provider) GC() {
	mp.sessionDao.deleteSessionByMaxLifeTime(mp.maxLifeTime)
}

// read session store by session id
func (mp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	sessionValue, err := mp.sessionDao.getSessionBySessionId(sessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		_, err := mp.sessionDao.insert(sessionId, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewMysqlStore(sessionId), nil
	}
	if len(sessionValue["contents"]) == 0 {
		return NewMysqlStore(sessionId), nil
	}

	data, err := mp.config.UnSerializeFunc(sessionValue["contents"])
	if err != nil {
		return nil, err
	}

	return NewMysqlStoreData(sessionId, data), nil
}

// regenerate session
func (mp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	sessionValue, err := mp.sessionDao.getSessionBySessionId(oldSessionId)
	if err != nil {
		return nil, err
	}
	if len(sessionValue) == 0 {
		// old sessionId not exists, insert new sessionId
		_, err := mp.sessionDao.insert(sessionId, "", time.Now().Unix())
		if err != nil {
			return nil, err
		}
		return NewMysqlStore(sessionId), nil
	}

	// delete old session
	_, err = mp.sessionDao.deleteBySessionId(oldSessionId)
	if err != nil {
		return nil, err
	}
	// insert new session
	_, err = mp.sessionDao.insert(sessionId, string(sessionValue["contents"]), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	return mp.ReadStore(sessionId)
}

// destroy session by sessionId
func (mp *Provider) Destroy(sessionId string) error {
	_, err := mp.sessionDao.deleteBySessionId(sessionId)
	return err
}

// session values count
func (mp *Provider) Count() int {
	return mp.sessionDao.countSessions()
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, provider)
}
