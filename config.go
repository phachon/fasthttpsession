package fasthttpsession

import "github.com/satori/go.uuid"

type Config struct {

	// cookie name
	CookieName string

	// cookie lifetime
	Lifetime int

	// is secure
	IsSecure bool

	// cookie expires time
	CookieExpires int

	// Domain
	Domain string

	// sessionId is in url query
	SIdIsInURLQuery bool

	// SessionId in url query name
	SIdInUrlQueryName string

	// sessionId is in http header
	SIdIsInHttpHeader bool

	// sessionId in http header name
	SIdInHttpHeaderName string

	// encrypt session data
	Encrypted bool

	// SessionIDGenerator should returns a random session id.
	// By default we will use a uuid impl package to generate
	// that, but developers can change that with simple assignment.
	SessionIdGeneratorFunc func() string

	// Encode the cookie value if not nil.
	// Should accept as first argument the cookie name (config.Name)
	//         as second argument the server's generated session id.
	// Should return the new session id, if error the session id setted to empty which is invalid.
	//
	// Note: Errors are not printed, so you have to know what you're doing,
	// and remember: if you use AES it only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	//
	// Defaults to nil
	EncodeFunc func(cookieName string, value interface{}) (string, error)

	// Decode the cookie value if not nil.
	// Should accept as first argument the cookie name (config.Name)
	//               as second second accepts the client's cookie value (the encoded session id).
	// Should return an error if decode operation failed.
	//
	// Note: Errors are not printed, so you have to know what you're doing,
	// and remember: if you use AES it only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	//
	// Defaults to nil
	DecodeFunc func(cookieName string, cookieValue string, v interface{}) error

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
	id, _ := uuid.NewV4()
	return id.String()
}

// encode cookie value
func (c *Config) Encode(cookieValue string) string {
	encode := c.EncodeFunc;
	if encode != nil {
		newVal, err := encode(c.CookieName, cookieValue)
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
	var cookieValueDecoded *string
	decode := c.DecodeFunc;
	if decode != nil {
		err := decode(c.CookieName, cookieValue, &cookieValueDecoded)
		if err == nil {
			cookieValue = *cookieValueDecoded
		} else {
			cookieValue = ""
		}
	}
	return cookieValue
}


type AdapterConfig struct {

}