package file

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/savsgio/fasthttpsession"
	"github.com/valyala/fasthttp"
)

// Store store struct
type Store struct {
	fasthttpsession.Store
}

// Save save store
func (fs *Store) Save(ctx *fasthttp.RequestCtx) error {

	fileProvider.lock.Lock()
	defer fileProvider.lock.Unlock()

	sessionID := fs.GetSessionID()

	_, _, fullFileName := fileProvider.getSessionFile(sessionID)

	if fileProvider.file.pathIsExists(fullFileName) {
		sessionMap := fs.GetAll()
		sessionInfo, _ := fileProvider.config.SerializeFunc(sessionMap)
		ioutil.WriteFile(fullFileName, sessionInfo, 0777)
		os.Chtimes(fullFileName, time.Now(), time.Now())
	}
	return nil
}
