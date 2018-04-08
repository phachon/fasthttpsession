package fasthttpsession

import (
	"fmt"
	"os"
	"github.com/valyala/fasthttp"
	"time"
)

type SessionAdapter interface {
	Init(*AdapterConfig) error
	Release(*fasthttp.RequestCtx) error
	SessionIdIsExit(string) bool
	GetSessionStore(string) SessionStore
}

type SessionStore interface {
	Get() string
	Set() error
	GetOnce() error
	Bind() error
	Destroy()
}

var adapters = make(map[string]SessionAdapter)

// register session adapter
func Register(adapterName string, adapter SessionAdapter)  {
	if adapters[adapterName] != nil {
		panic("session: session adapter "+ adapterName +" already registered!")
	}
	if adapter == nil {
		panic("session: session adapter "+ adapterName +" is nil!")
	}

	adapters[adapterName] = adapter
}

// Session
type Session struct {
	Adapter SessionAdapter
	Config *Config
}

// return new Session, default adapter file
func NewSession(config *Config) *Session {
	session := &Session{
		Config: config,
	}

	// default adapter file
	session.SetAdapter("file", &AdapterConfig{})

	return session
}

// set session adapter and adapter config
func (s *Session) SetAdapter(adapterName string, adapterConfig *AdapterConfig) error {
	adapter, ok := adapters[adapterName]
	if !ok {
		printError("session: session adapter "+adapterName+" not found!")
	}

	err := adapter.Init(adapterConfig)
	if err != nil {
		return err
	}

	s.Adapter = adapter
	return nil
}

// session start
// return session store
func (s *Session) Start(ctx *fasthttp.RequestCtx) SessionStore {

	sId := s.GetSessionId(ctx)

	// if sessionId is not empty, check is exit in session adapter
	if sId != "" && s.Adapter.SessionIdIsExit(sId) {
		return s.Adapter.GetSessionStore(sId)
	}
	// new session id
	sId = s.Config.SessionIdGeneratorFunc()
	if sId == "" {
		printError("sessionId generator is empty")
	}

	// add cookie
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(s.Config.CookieName)
	cookie.SetValue(sId)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetDomain(s.Config.Domain)
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	if s.Config.CookieExpires >= 0 {
		if s.Config.CookieExpires == 0 {
			// = 0 unlimited life
			cookie.SetExpire(fasthttp.CookieExpireUnlimited)
		} else {
			// > 0
			cookie.SetExpire(time.Now().Add(s.Config.CookieExpires))
		}
	}

	if ctx.IsTLS() && s.Config.IsSecure {
		cookie.SetSecure(true)
	}

	// encode cookie value
	cookie.SetValue(s.Config.Encode(string(cookie.Value())))
	ctx.Response.Header.SetCookie(cookie)

	return sId
}

func (s *Session) Release(ctx *fasthttp.RequestCtx) {

}

// get fasthttp session id
// 1. get session id by reading from cookie
// 2. if not exist, get session id from query or http headers
func (s *Session) GetSessionId(ctx *fasthttp.RequestCtx) (string) {

	cookieByte := ctx.Request.Header.Cookie(s.Config.CookieName)
	if len(cookieByte) > 0 {
		return (string(cookieByte))
	}

	if s.Config.SIdIsInURLQuery {
		cookieFormValue := ctx.FormValue(s.Config.SIdInUrlQueryName)
		if len(cookieFormValue) > 0 {
			return string(cookieFormValue)
		}
	}

	if s.Config.SIdIsInHttpHeader {
		cookieHeader := ctx.Request.Header.Peek(s.Config.SIdInHttpHeaderName)
		if len(cookieHeader) > 0 {
			return string(cookieHeader)
		}
	}

	return ""
}

// Set cookie with https.
func (s *Session) isSecure(ctx *fasthttp.RequestCtx) bool {
	if !s.Config.IsSecure {
		return false
	}
	if string(ctx.Request.URI().Scheme()) != "" {
		return string(ctx.Request.URI().Scheme()) == "https"
	}
	if !ctx.IsTLS() {
		return false
	}
	return true
}

func printError(msg string)  {
	fmt.Println(msg)
	os.Exit(100)
}