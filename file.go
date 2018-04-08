package fasthttpsession

import "github.com/valyala/fasthttp"

// file session adapter

const Session_Adapter_File = "file"

func NewAdapterFile() *AdapterFile {
	return &AdapterFile{}
}

type AdapterFile struct {

}

func (f *AdapterFile) Init(config *AdapterConfig) error {
	return nil
}

func (f *AdapterFile) Release(ctx *fasthttp.RequestCtx) {

}

func (f *AdapterFile) SessionIsExit(sessionId string) bool {

}

func init()  {
	Register(Session_Adapter_File, NewAdapterFile())
}