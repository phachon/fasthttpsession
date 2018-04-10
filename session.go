package fasthttpsession

import (
	"github.com/valyala/fasthttp"
	"errors"
)

// Session struct
type Session struct {
	provider Provider
	config *Config
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


// return new Session, default provider file
func NewSession(cfg *Config) *Session {
	session := &Session{
		config: cfg,
	}

	return session
}

// set session provider and provider config
func (s *Session) SetProvider(providerName string, providerConfig ProviderConfig) error {
	provider, ok := providers[providerName]
	if !ok {
		return errors.New("session: session provider "+providerName+" not found!")
	}

	err := provider.Init(providerConfig)
	if err != nil {
		return err
	}

	s.provider = provider
	return nil
}

// session start
// return session provider
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
	// set cookie
	NewCookie().Set(ctx,
		s.config.CookieName,
		encodeCookieValue,
		s.config.Domain,
		s.config.Expires,
		s.config.Secure)
	return
}

// get session id
// 1. get session id by reading from cookie
// 2. get session id from query or http headers
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

// destroy session by its id in fasthttp request
func (s *Session) Destroy(ctx *fasthttp.RequestCtx) {

	// delete header if sessionId in http Header
	if s.config.SessionIdInHttpHeader {
		ctx.Request.Header.Del(s.config.SessionNameInHttpHeader)
		ctx.Response.Header.Del(s.config.SessionNameInHttpHeader)
	}

	cookieValue := string(ctx.Request.Header.Cookie(s.config.CookieName))
	if cookieValue == "" {
		return
	}

	//sessionId, _ := url.QueryUnescape(string(cookieByte))
	sessionId, _ := s.config.Decode(cookieValue)
	s.provider.Destroy(sessionId)

	// delete cookie by cookieName
	NewCookie().Delete(ctx, s.config.CookieName)
}