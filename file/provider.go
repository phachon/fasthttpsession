package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/savsgio/fasthttpsession"
)

// session file provider

const ProviderName = "file"

var (
	fileProvider = NewProvider()
	encrypt      = fasthttpsession.NewEncrypt()
)

type Provider struct {
	lock        sync.RWMutex
	file        *file
	config      *Config
	maxLifeTime int64
}

// new file provider
func NewProvider() *Provider {
	return &Provider{
		file:   &file{},
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
func (fp *Provider) GC() {

	files, err := fp.file.walkDir(fp.config.SavePath, fp.config.Suffix)
	if err == nil {
		for _, file := range files {
			if time.Now().Unix() >= (fp.maxLifeTime + fp.file.getModifyTime(file)) {
				fp.lock.Lock()
				filename := filepath.Base(file)
				sessionID := strings.TrimRight(filename, fp.config.Suffix)
				fp.removeSessionFile(sessionID)
				fp.lock.Unlock()
			}
		}
	}
}

// read session store by session id
func (fp *Provider) ReadStore(sessionID string) (fasthttpsession.SessionStore, error) {

	fp.lock.Lock()
	defer fp.lock.Unlock()
	store := &Store{}

	filePath, _, fullFileName := fp.getSessionFile(sessionID)

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
		store.Init(sessionID, value)

		return store, nil
	}

	os.MkdirAll(filePath, 0777)

	err := fp.file.createFile(fullFileName)
	if err != nil {
		return store, err
	}
	store.Init(sessionID, map[string]interface{}{})

	return store, nil
}

// regenerate session
func (fp *Provider) Regenerate(oldSessionId string, sessionID string) (fasthttpsession.SessionStore, error) {

	fp.lock.Lock()
	defer fp.lock.Unlock()
	store := &Store{}

	_, _, oldFullFileName := fp.getSessionFile(oldSessionId)
	filePath, _, fullFileName := fp.getSessionFile(sessionID)

	if fp.file.pathIsExists(fullFileName) {
		return store, errors.New("new sessionID file exist")
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
		store.Init(sessionID, value)

		return store, nil
	}

	store.Init(sessionID, map[string]interface{}{})

	return store, nil
}

// destroy session by sessionID
func (fp *Provider) Destroy(sessionID string) error {

	fp.lock.Lock()
	defer fp.lock.Unlock()

	_, _, fullFileName := fp.getSessionFile(sessionID)
	if fp.file.pathIsExists(fullFileName) {
		fp.removeSessionFile(sessionID)
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
func (fp *Provider) getSessionFile(sessionID string) (string, string, string) {
	filePath := path.Join(fp.config.SavePath, string(sessionID[0]), string(sessionID[1]))
	filename := sessionID + fp.config.Suffix
	fullFilename := filepath.Join(filePath, filename)

	return filePath, filename, fullFilename
}

// remove session file
func (fp *Provider) removeSessionFile(sessionID string) {

	filePath, _, fullFileName := fp.getSessionFile(sessionID)
	os.Remove(fullFileName)

	// remove empty dir
	s, _ := ioutil.ReadDir(filePath)
	if len(s) == 0 {
		os.RemoveAll(filePath)
	}
	filePath1 := path.Join(fp.config.SavePath, string(sessionID[0]))
	s, _ = ioutil.ReadDir(filePath1)
	if len(s) == 0 {
		os.RemoveAll(filePath1)
	}
}

// register session provider
func init() {
	fasthttpsession.Register(ProviderName, fileProvider)
}
