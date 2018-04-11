package fasthttpsession

import (
	"github.com/satori/go.uuid"
	"time"
)

// new default config
func NewDefaultConfig() *Config {
	config := &Config{
		CookieName: "fasthttpsessionid",
		Domain: "",
		Expires: time.Hour * 2,
		GCLifetime: 5,
		SessionLifetime: 60,
		Secure: true,
		SessionIdInURLQuery: false,
		SessionNameInUrlQuery: "",
		SessionIdInHttpHeader: false,
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

	// gc life time
	GCLifetime int64

	// session life time
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
	SessionIdGeneratorFunc func() (string, error)

	// Encode the cookie value if not nil.
	EncodeFunc func(cookieValue string) (string, error)

	// Decode the cookie value if not nil.
	DecodeFunc func(cookieValue string) (string, error)

}

type ProviderConfig interface {
	Name() string
}

// sessionId generator
func (c *Config) SessionIdGenerator() (string, error) {
	sessionIdGenerator := c.SessionIdGeneratorFunc
	if sessionIdGenerator == nil {
		return c.defaultSessionIdGenerator()
	}

	return sessionIdGenerator()
}

// default sessionId generator => uuid
func (c *Config) defaultSessionIdGenerator() (string, error) {
	id, err := uuid.NewV4()
	return id.String(), err
}

// encode cookie value
func (c *Config) Encode(cookieValue string) (string, error) {
	encode := c.EncodeFunc;
	if encode != nil {
		newVal, err := encode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue, nil
}

// decode cookie value
func (c *Config) Decode(cookieValue string) (string, error) {
	if cookieValue == "" {
		return "", nil
	}
	decode := c.DecodeFunc;
	if decode != nil {
		newVal, err := decode(cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}
	return cookieValue, nil
}