package file

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"time"
)

type Store struct {
	fasthttpsession.Store
}

// save store
func (fs *Store) Save(ctx *fasthttp.RequestCtx) error {

	fileProvider.lock.Lock()
	defer fileProvider.lock.Unlock()

	sessionId := fs.GetSessionId()

	_, _, fullFileName := fileProvider.getSessionFile(sessionId)

	if fileProvider.file.pathIsExists(fullFileName) {
		sessionMap := fs.GetAll()
		sessionInfo, _ := fileProvider.config.SerializeFunc(sessionMap)
		ioutil.WriteFile(fullFileName, sessionInfo, 0777)
		os.Chtimes(fullFileName, time.Now(), time.Now())
	}
	return nil
}
