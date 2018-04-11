package fasthttpsession

import (
	"github.com/valyala/fasthttp"
	"errors"
	"time"
	"fmt"
)

// Session struct
type Session struct {
	provider Provider
	config  *Config
	cookie  *Cookie
}

var providers = make(map[string]Provider)

// register session provider
func Register(providerName string, provider Provider)  {
	if providers[providerName] != nil {
		panic("session: session provider "+ providerName +" already registered!")
	}
	if provider == nil {
		panic("session: session provider "+ providerName +" is nil!")
	}

	providers[providerName] = provider
}

// return new Session
func NewSession(cfg *Config) *Session {
	session := &Session{
		config: cfg,
		cookie: NewCookie(),
	}

	return session
}

// set session provider and provider config
func (s *Session) SetProvider(providerName string, providerConfig ProviderConfig) error {
	provider, ok := providers[providerName]
	if !ok {
		return errors.New("session: session provider "+providerName+" not registered!")
	}
	err := provider.Init(providerConfig)
	if err != nil {
		return err
	}
	s.provider = provider

	// start gc
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				panic(errors.New(fmt.Sprintf("session gc crash, %v", e)))
			}
		}()
		s.gc()
	}()
	return nil
}

// start session gc process.
func (s *Session) gc() {
	for {
		select {
		case <-time.After(time.Duration(s.config.GCLifetime) * time.Second):
			s.provider.GC(s.config.SessionLifetime)
		}
	}
}

// session start
// 1. get sessionId from fasthttp ctx
// 2. if sessionId is empty, generator sessionId and set response Set-Cookie
// 3. return session provider store
func (s *Session) Start(ctx *fasthttp.RequestCtx) (sessionStore SessionStore, err error) {

	sessionId, err := s.GetSessionId(ctx)
	if err != nil {
		return
	}
	// if sessionId is not empty, check is exit in session provider
	if sessionId != "" && s.provider.SessionIdIsExist(sessionId) {
		return s.provider.ReadStore(sessionId)
	}
	// new session id
	sessionId, err = s.config.SessionIdGenerator()
	if err != nil{
		return
	}
	if sessionId == "" {
		return sessionStore, errors.New("generator sessionId  is empty")
	}
	sessionStore, err = s.provider.ReadStore(sessionId)
	if err != nil {
		return
	}

	// encode cookie value
	encodeCookieValue, err := s.config.Encode(sessionId)
	if err != nil {
		return
	}
	// set response cookie
	s.cookie.Set(ctx,
		s.config.CookieName,
		encodeCookieValue,
		s.config.Domain,
		s.config.Expires,
		s.config.Secure)

	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Set(s.config.SessionNameInHttpHeader, sessionId)
		ctx.Response.Header.Set(s.config.SessionNameInHttpHeader, sessionId)
	}

	return
}

// get session id
// 1. get session id by reading from cookie
// 2. get session id from query
// 3. get session id from http headers
func (s *Session) GetSessionId(ctx *fasthttp.RequestCtx) (string, error) {

	cookieByte := ctx.Request.Header.Cookie(s.config.CookieName)
	if len(cookieByte) > 0 {
		return s.config.Decode(string(cookieByte))
	}

	if s.config.SessionIdInURLQuery {
		cookieFormValue := ctx.FormValue(s.config.SessionNameInUrlQuery)
		if len(cookieFormValue) > 0 {
			return s.config.Decode(string(cookieFormValue))
		}
	}

	if s.config.SessionIdInHttpHeader {
		cookieHeader := ctx.Request.Header.Peek(s.config.SessionNameInHttpHeader)
		if len(cookieHeader) > 0 {
			return s.config.Decode(string(cookieHeader))
		}
	}

	return "", nil
}

// regenerate a session id for this SessionStore
func (s *Session) Regenerate(ctx *fasthttp.RequestCtx) (sessionStore SessionStore, err error) {

	// generator new session id
	sessionId, err := s.config.SessionIdGenerator()
	if err != nil{
		return
	}
	if sessionId == "" {
		return sessionStore, errors.New("generator sessionId  is empty")
	}
	// encode cookie value
	encodeCookieValue, err := s.config.Encode(sessionId)
	if err != nil {
		return
	}

	oldSessionId, err := s.GetSessionId(ctx)
	if err != nil {
		return
	}
	// regenerate provider session store
	if oldSessionId != "" {
		sessionStore, err = s.provider.Regenerate(oldSessionId, sessionId)
	}else {
		sessionStore, err = s.provider.ReadStore(sessionId)
	}
	if err != nil {
		return
	}

	// reset response cookie
	s.cookie.Set(ctx,
		s.config.CookieName,
		encodeCookieValue,
		s.config.Domain,
		s.config.Expires,
		s.config.Secure)

	// reset http header
	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Set(s.config.SessionNameInHttpHeader, sessionId)
		ctx.Response.Header.Set(s.config.SessionNameInHttpHeader, sessionId)
	}

	return
}

// destroy session in fasthttp ctx
func (s *Session) Destroy(ctx *fasthttp.RequestCtx) {

	// delete header if sessionId in http Header
	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHttpHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHttpHeader)
	}

	cookieValue := s.cookie.Get(ctx, s.config.CookieName)
	if cookieValue == "" {
		return
	}

	sessionId, err := s.config.Decode(cookieValue)
	if err != nil {
		return
	}
	s.provider.Destroy(sessionId)

	// delete cookie by cookieName
	s.cookie.Delete(ctx, s.config.CookieName)
}