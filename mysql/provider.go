package mysql

import (
	"github.com/phachon/fasthttpsession"
	"errors"
	"reflect"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
)

// session mysql provider

// session Table structure
// #-- ----------------------------------------------------------
// #-- session table
// #-- ----------------------------------------------------------
// DROP TABLE IF EXISTS `session`;
// CREATE TABLE `session` (
//   `session_id` varchar(24) NOT NULL DEFAULT '' COMMENT 'Session id',
//   `last_active` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Last active time',
//   `contents` varchar(1000) NOT NULL DEFAULT '' COMMENT 'Session data',
//   PRIMARY KEY (`session_id`),
//   KEY `last_active` (`last_active`)
// ) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='session table';
//

const ProviderName = "mysql"

var (
	utils = fasthttpsession.NewUtils()
	provider = NewProvider()
)

type Provider struct {
	config *Config
	values *fasthttpsession.CCMap
	mysqlConn *sql.DB
	maxLifeTime int64
}

// new mysql provider
func NewProvider() *Provider {
	return &Provider{
		config: &Config{},
		values: fasthttpsession.NewDefaultCCMap(),
		mysqlConn: &sql.DB{},
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

	if mp.config.Host == "" {
		return errors.New("session mysql provider init error, config Host not empty")
	}
	if mp.config.Port == 0 {
		return errors.New("session mysql provider init error, config Port not empty")
	}

	// init mysql conn
	mysqlConn, err := getMysqlConn(mp.config.getMysqlDSN())
	if err != nil {
		return err
	}
	mysqlConn.SetMaxOpenConns(mp.config.SetMaxIdleConn)
	mysqlConn.SetMaxIdleConns(mp.config.SetMaxIdleConn)

	return mysqlConn.Ping()
}

// not need gc
func (mp *Provider) NeedGC() bool {
	return false
}

// session mysql provider not need garbage collection
func (mp *Provider) GC(sessionLifetime int64) {}


// read session store by session id
func (mp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	sqlStr := fmt.Sprintf("SELECT * FROM %s where session_id=?", mp.config.TableName)

	res, err := getRow(mp.mysqlConn, sqlStr, sessionId)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		//todo insert sessionId
		return NewmysqlStore(sessionId), nil
	}
	if len(res["contents"]) == 0 {
		return NewmysqlStore(sessionId), nil
	}

	data, err := utils.GobDecode(res["content"])
	if err != nil {
		return nil, err
	}

	return NewmysqlStoreData(sessionId, data), nil
}

// regenerate session
func (mp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {


	sqlStr := fmt.Sprintf("SELECT * FROM %s where session_id=?", mp.config.TableName)

	res, err := getRow(mp.mysqlConn, sqlStr, oldSessionId)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		// old sessionId not exists

		insertSql := fmt.Sprintf("INSERT into ")
		execute(db, )


		return NewmysqlStore(sessionId), nil
	}
	if len(res["contents"]) == 0 {
		return NewmysqlStore(sessionId), nil
	}

	data, err := utils.GobDecode(res["content"])
	if err != nil {
		return nil, err
	}

	return NewmysqlStoreData(sessionId, data), nil

	conn := rp.mysqlPool.Get()
	defer conn.Close()

	existed, err := mysql.Int(conn.Do("EXISTS", rp.getmysqlSessionKey(oldSessionId)))
	if err != nil || existed == 0 {
		// false
		conn.Do("SET", rp.getmysqlSessionKey(sessionId), "", "EX", rp.maxLifeTime)
		return NewmysqlStore(sessionId), nil
	}
	// true
	conn.Do("RENAME", rp.getmysqlSessionKey(oldSessionId), rp.getmysqlSessionKey(sessionId))
	conn.Do("EXPIRE", rp.getmysqlSessionKey(sessionId), rp.maxLifeTime)

	return rp.ReadStore(sessionId)
}

// destroy session by sessionId
func (rp *Provider) Destroy(sessionId string) error {
	conn := rp.mysqlPool.Get()
	defer conn.Close()

	existed, err := mysql.Int(conn.Do("EXISTS", rp.getmysqlSessionKey(sessionId)))
	if err != nil || existed == 0 {
		return nil
	}
	conn.Do("DEL", rp.getmysqlSessionKey(sessionId))
	return nil
}

// session values count
func (rp *Provider) Count() int {
	conn := rp.mysqlPool.Get()
	defer conn.Close()

	replyMap, err := mysql.Strings(conn.Do("KEYS", rp.config.KeyPrefix+":*"))
	if err != nil {
		return 0
	}
	return len(replyMap)
}

// get mysql session key, prefix:sessionId
func (rp *Provider) getmysqlSessionKey(sessionId string) string {
	return rp.config.KeyPrefix+":"+sessionId
}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, provider)
}