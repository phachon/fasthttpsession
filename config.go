package fasthttpsession

import (
	"time"

	"github.com/satori/go.uuid"
)

var (
	defaultCookieName = "fasthttpsessionid"

	defaultExpires = time.Hour * 2

	defaultGCLifetime = int64(3)
)

// new default config
func NewDefaultConfig() *Config {
	config := &Config{
		CookieName:              defaultCookieName,
		Domain:                  "",
		Expires:                 defaultExpires,
		GCLifetime:              defaultGCLifetime,
		SessionLifetime:         60,
		Secure:                  true,
		SessionIdInURLQuery:     false,
		SessionNameInUrlQuery:   "",
		SessionIdInHttpHeader:   false,
		SessionNameInHttpHeader: "",
	}

	// default sessionIdGeneratorFunc
	config.SessionIdGeneratorFunc = config.defaultSessionIdGenerator

	return config
}

type Config struct {

	// cookie name
	CookieName string

	// cookie domain
	Domain string

	// If you want to delete the cookie when the browser closes, set it to -1.
	//
	//  0 means no expire, (24 years)
	// -1 means when browser closes
	// >0 is the time.Duration which the session cookies should expire.
	Expires time.Duration

	// gc life time(s)
	GCLifetime int64

	// session life time(s)
	SessionLifetime int64

	// set whether to pass this bar cookie only through HTTPS
	Secure bool

	// sessionId is in url query
	SessionIdInURLQuery bool

	// sessionName in url query
	SessionNameInUrlQuery string

	// sessionId is in http header
	SessionIdInHttpHeader bool

	// sessionName in http header
	SessionNameInHttpHeader string

	// SessionIdGeneratorFunc should returns a random session id.
	SessionIdGeneratorFunc func() string

	// Encode the cookie value if not nil.
	EncodeFunc func(cookieValue string) (string, error)

	// Decode the cookie value if not nil.
	DecodeFunc func(cookieValue string) (string, error)
}

// sessionId generator
func (c *Config) SessionIdGenerator() string {
	sessionIdGenerator := c.SessionIdGeneratorFunc
	if sessionIdGenerator == nil {
		return c.defaultSessionIdGenerator()
	}

	return sessionIdGenerator()
}

// default sessionId generator => uuid
func (c *Config) defaultSessionIdGenerator() string {
	id, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return id.String()
}

// encode cookie value
func (c *Config) Encode(cookieValue string) string {
	encode := c.EncodeFunc
	if encode != nil {
		newVal, err := encode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue
}

// decode cookie value
func (c *Config) Decode(cookieValue string) string {
	if cookieValue == "" {
		return ""
	}
	decode := c.DecodeFunc
	if decode != nil {
		newVal, err := decode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue
}
