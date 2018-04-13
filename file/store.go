package file

import (
	"github.com/phachon/fasthttpsession"
	"github.com/valyala/fasthttp"
)

type Store struct {
	fasthttpsession.Store
}

// save store
func (ms *Store) Save(ctx *fasthttp.RequestCtx) error {
	ms.Lock.Lock()
	defer ms.Lock.Unlock()

	return nil
}