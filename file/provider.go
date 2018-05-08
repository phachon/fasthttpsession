package file

import (
	"github.com/phachon/fasthttpsession"
	"errors"
	"reflect"
	"os"
	"path"
	"time"
	"sync"
	"path/filepath"
	"io/ioutil"
	"strings"
)

// session file provider

const ProviderName = "file"

var (
	fileProvider = NewProvider()
	encrypt = fasthttpsession.NewEncrypt()
)

type Provider struct {
	lock sync.RWMutex
	file *file
	config *Config
	maxLifeTime int64
}

// new file provider
func NewProvider() *Provider {
	return &Provider{
		file: &file{},
		config: &Config{},
	}
}

// init provider config
func (fp *Provider) Init(lifeTime int64, fileConfig fasthttpsession.ProviderConfig) error {
	if fileConfig.Name() != ProviderName {
		return errors.New("session file provider init error, config must file config")
	}

	vc := reflect.ValueOf(fileConfig)
	fc := vc.Interface().(*Config)
	fp.config = fc

	if fp.config.SavePath == "" {
		return errors.New("session file provider init error, config savePath not empty")
	}
	if fp.config.SerializeFunc == nil {
		fp.config.SerializeFunc = encrypt.GobEncode
	}
	if fp.config.UnSerializeFunc == nil {
		fp.config.UnSerializeFunc = encrypt.GobDecode
	}

	fp.maxLifeTime = lifeTime

	// create save path
	os.MkdirAll(fp.config.SavePath, 0777)

	return nil
}

// need gc
func (fp *Provider) NeedGC() bool {
	return true
}

// session garbage collection
func (fp *Provider) GC(sessionLifetime int64) {

	files, err := fp.file.walkDir(fp.config.SavePath, fp.config.Suffix)
	if err == nil {
		for _, file := range files {
			if time.Now().Unix() >= (sessionLifetime + fp.file.getModifyTime(file)) {
				fp.lock.Lock()
				filename := filepath.Base(file)
				sessionId := strings.TrimRight(filename, fp.config.Suffix)
				fp.removeSessionFile(sessionId)
				fp.lock.Unlock()
			}
		}
	}
}

// read session store by session id
func (fp *Provider) ReadStore(sessionId string) (fasthttpsession.SessionStore, error) {

	fp.lock.Lock()
	defer fp.lock.Unlock()
	store := &Store{}

	filePath, _, fullFileName := fp.getSessionFile(sessionId)

	// file is exist
	if fp.file.pathIsExists(fullFileName) {
		sessionInfo, err := fp.file.getContent(fullFileName)
		if err != nil {
			return store, err
		}

		// unserialize sessionInfo
		value, err := fp.config.UnSerializeFunc(sessionInfo)
		if err != nil {
			return store, err
		}
		store.Init(sessionId, value)

		return store, nil
	}

	os.MkdirAll(filePath, 0777)

	err := fp.file.createFile(fullFileName)
	if err != nil {
		return store, err
	}
	store.Init(sessionId, map[string]interface{}{})

	return store, nil
}


// regenerate session
func (fp *Provider) Regenerate(oldSessionId string, sessionId string) (fasthttpsession.SessionStore, error) {

	fp.lock.Lock()
	defer fp.lock.Unlock()
	store := &Store{}

	_, _, oldFullFileName := fp.getSessionFile(oldSessionId)
	filePath, _, fullFileName := fp.getSessionFile(sessionId)

	if fp.file.pathIsExists(fullFileName) {
		return store, errors.New("new sessionId file exist")
	}
	// create new session file
	os.MkdirAll(filePath, 0777)
	err := fp.file.createFile(fullFileName)
	if err != nil {
		return store, err
	}

	if fp.file.pathIsExists(oldFullFileName) {
		// read old session info
		sessionInfo, err := fp.file.getContent(fullFileName)
		if err != nil {
			return store, err
		}
		// write new session file
		ioutil.WriteFile(fullFileName, sessionInfo, 0777)
		// remove old session file
		fp.removeSessionFile(oldSessionId)
		// update new session file time
		os.Chtimes(fullFileName, time.Now(), time.Now())

		// unserialize sessionInfo
		value, err := fp.config.UnSerializeFunc(sessionInfo)
		if err != nil {
			return store, err
		}
		store.Init(sessionId, value)

		return store, nil
	}

	store.Init(sessionId, map[string]interface{}{})

	return store, nil
}

// destroy session by sessionId
func (fp *Provider) Destroy(sessionId string) error {

	fp.lock.Lock()
	defer fp.lock.Unlock()

	_, _, fullFileName  := fp.getSessionFile(sessionId)
	if fp.file.pathIsExists(fullFileName) {
		fp.removeSessionFile(sessionId)
	}

	return nil
}

// session values count
func (fp *Provider) Count() int {
	fp.lock.Lock()
	defer fp.lock.Unlock()

	count, _ := fp.file.count(fp.config.SavePath, fp.config.Suffix)

	return count
}

// get session filePath, filename, fullFilename
func (fp *Provider) getSessionFile(sessionId string) (string, string, string) {
	filePath := path.Join(fp.config.SavePath, string(sessionId[0]), string(sessionId[1]))
	filename := sessionId + fp.config.Suffix
	fullFilename := filepath.Join(filePath, filename)

	return filePath, filename, fullFilename
}

// remove session file
func (fp *Provider) removeSessionFile(sessionId string) {

	filePath, _, fullFileName  := fp.getSessionFile(sessionId)
	os.Remove(fullFileName)

	// remove empty dir
	s, _ := ioutil.ReadDir(filePath)
	if len(s) == 0 {
		os.RemoveAll(filePath)
	}
	filePath1 := path.Join(fp.config.SavePath, string(sessionId[0]))
	s, _ = ioutil.ReadDir(filePath1)
	if len(s) == 0 {
		os.RemoveAll(filePath1)
	}
}

// register session provider
func init()  {
	fasthttpsession.Register(ProviderName, fileProvider)
}