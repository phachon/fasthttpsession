package fasthttpsession

import (
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var version = "v0.0.1"

// Session session struct
type Session struct {
	provider Provider
	config   *Config
	cookie   *Cookie
}

var providers = make(map[string]Provider)

// Register register session provider
func Register(providerName string, provider Provider) {
	if providers[providerName] != nil {
		panic("session register error, provider " + providerName + " already registered!")
	}
	if provider == nil {
		panic("session register error, provider " + providerName + " is nil!")
	}

	providers[providerName] = provider
}

// NewSession return new Session
func NewSession(cfg *Config) *Session {

	if cfg.CookieName == "" {
		cfg.CookieName = defaultCookieName
	}
	if cfg.GCLifetime == 0 {
		cfg.GCLifetime = defaultGCLifetime
	}
	if cfg.SessionLifetime == 0 {
		cfg.SessionLifetime = cfg.GCLifetime
	}
	if cfg.SessionIdGeneratorFunc == nil {
		cfg.SessionIdGeneratorFunc = cfg.defaultSessionIdGenerator
	}

	session := &Session{
		config: cfg,
		cookie: NewCookie(),
	}

	return session
}

// SetProvider set session provider and provider config
func (s *Session) SetProvider(providerName string, providerConfig ProviderConfig) error {
	provider, ok := providers[providerName]
	if !ok {
		return errors.New("session set provider error, " + providerName + " not registered!")
	}
	err := provider.Init(s.config.SessionLifetime, providerConfig)
	if err != nil {
		return err
	}
	s.provider = provider

	// start gc
	if s.provider.NeedGC() {
		go func() {
			defer func() {
				e := recover()
				if e != nil {
					panic(fmt.Errorf("session gc crash, %v", e))
				}
			}()
			s.gc()
		}()
	}
	return nil
}

// start session gc process.
func (s *Session) gc() {
	for {
		select {
		case <-time.After(time.Duration(s.config.GCLifetime) * time.Second):
			s.provider.GC()
		}
	}
}

// Start session start
// 1. get sessionID from fasthttp ctx
// 2. if sessionID is empty, generator sessionID and set response Set-Cookie
// 3. return session provider store
func (s *Session) Start(ctx *fasthttp.RequestCtx) (sessionStore SessionStore, err error) {

	if s.provider == nil {
		return sessionStore, errors.New("session start error, not set provider")
	}

	sessionID := s.GetSessionId(ctx)
	if sessionID == "" {
		// new generator session id
		sessionID = s.config.SessionIdGenerator()
		if sessionID == "" {
			return sessionStore, errors.New("session generator sessionID is empty")
		}
	}
	// read provider session store
	sessionStore, err = s.provider.ReadStore(sessionID)
	if err != nil {
		return
	}

	// encode cookie value
	encodeCookieValue := s.config.Encode(sessionID)

	// set response cookie
	s.cookie.Set(ctx,
		s.config.CookieName,
		encodeCookieValue,
		s.config.Domain,
		s.config.Expires,
		s.config.Secure)

	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Set(s.config.SessionNameInHttpHeader, sessionID)
		ctx.Response.Header.Set(s.config.SessionNameInHttpHeader, sessionID)
	}

	return
}

// GetSessionId get session id
// 1. get session id by reading from cookie
// 2. get session id from query
// 3. get session id from http headers
func (s *Session) GetSessionId(ctx *fasthttp.RequestCtx) string {

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

	return ""
}

// Regenerate regenerate a session id for this SessionStore
func (s *Session) Regenerate(ctx *fasthttp.RequestCtx) (sessionStore SessionStore, err error) {

	if s.provider == nil {
		return sessionStore, errors.New("session regenerate error, not set provider")
	}

	// generator new session id
	sessionID := s.config.SessionIdGenerator()
	if sessionID == "" {
		return sessionStore, errors.New("session generator sessionID is empty")
	}
	// encode cookie value
	encodeCookieValue := s.config.Encode(sessionID)

	oldSessionId := s.GetSessionId(ctx)
	// regenerate provider session store
	if oldSessionId != "" {
		sessionStore, err = s.provider.Regenerate(oldSessionId, sessionID)
	} else {
		sessionStore, err = s.provider.ReadStore(sessionID)
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
		ctx.Request.Header.Set(s.config.SessionNameInHttpHeader, sessionID)
		ctx.Response.Header.Set(s.config.SessionNameInHttpHeader, sessionID)
	}

	return
}

// Destroy destroy session in fasthttp ctx
func (s *Session) Destroy(ctx *fasthttp.RequestCtx) {

	// delete header if sessionID in http Header
	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHttpHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHttpHeader)
	}

	cookieValue := s.cookie.Get(ctx, s.config.CookieName)
	if cookieValue == "" {
		return
	}

	sessionID := s.config.Decode(cookieValue)
	s.provider.Destroy(sessionID)

	// delete cookie by cookieName
	s.cookie.Delete(ctx, s.config.CookieName)
}

// Version return fasthttpsession version
func Version() string {
	return version
}
