package fasthttpsession

import (
	"github.com/valyala/fasthttp"
	"strings"
	"time"
	"strconv"
)

type Cookie struct  {

}

// Get cookie
func (c *Cookie) GetCookie(ctx *fasthttp.RequestCtx, name string) (value string) {
	cookieByte := ctx.Request.Header.Cookie(name)
	if cookieByte != nil {
		value = string(cookieByte)
	}
	return
}

func (c *Cookie) RequestAddCookie(ctx *fasthttp.RequestCtx, cookie *fasthttp.Cookie)  {

}

// Add cookie
func (c *Cookie) SetResponseCookie(ctx *fasthttp.RequestCtx, cookie *fasthttp.Cookie) {
	ctx.Response.Header.SetCookie(cookie)
}

// Remove Cookie
func (c *Cookie) RemoveCookie(ctx *fasthttp.RequestCtx, cookieName string) {
	ctx.Response.Header.DelCookie(cookieName)

	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(cookieName)
	cookie.SetValue("")
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	exp := time.Now().Add(-time.Duration(1) * time.Minute) //RFC says 1 second, but let's do it 1 minute to make sure is working...
	cookie.SetExpire(exp)
	c.SetResponseCookie(ctx, cookie)
	fasthttp.ReleaseCookie(cookie)
	// delete request's cookie also, which is temporary available
	ctx.Request.Header.DelCookie(cookieName)
}

// IsValidCookieDomain returns true if the receiver is a valid domain to set
// valid means that is recognised as 'domain' by the browser, so it(the cookie) can be shared with subdomains also
func (c *Cookie) IsValidCookieDomain(domain string) bool {
	if domain == "0.0.0.0" || domain == "127.0.0.1" {
		// for these type of hosts, we can't allow subdomains persistence,
		// the web browser doesn't understand the mysubdomain.0.0.0.0 and mysubdomain.127.0.0.1 mysubdomain.32.196.56.181. as scorrectly ubdomains because of the many dots
		// so don't set a cookie domain here, let browser handle this
		return false
	}

	dotLen := strings.Count(domain, ".")
	if dotLen == 0 {
		// we don't have a domain, maybe something like 'localhost', browser doesn't see the .localhost as wildcard subdomain+domain
		return false
	}
	if dotLen >= 3 {
		if lastDotIdx := strings.LastIndexByte(domain, '.'); lastDotIdx != -1 {
			// chekc the last part, if it's number then propably it's ip
			if len(domain) > lastDotIdx+1 {
				_, err := strconv.Atoi(domain[lastDotIdx+1:])
				if err == nil {
					return false
				}
			}
		}
	}

	return true
}
