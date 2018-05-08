[![logo](./logo.png)](https://github.com/phachon/fasthttpsession)

[![build](https://img.shields.io/shippable/5444c5ecb904a4b21567b0ff.svg)](https://travis-ci.org/phachon/fasthttpsession)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/phachon/fasthttpsession)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phachon/fasthttpsession/master/LICENSE)
[![go_Report](https://goreportcard.com/badge/github.com/phachon/fasthttpsession)](https://goreportcard.com/report/github.com/phachon/fasthttpsession)
[![release](https://img.shields.io/github/release/phachon/fasthttpsession.svg?style=flat)](https://github.com/phachon/fasthttpsession/releases) 
[![powered_by](https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat)]()
[![platforms](https://img.shields.io/badge/platform-All-yellow.svg?style=flat)]()

fasthttpsession is a fast and powerful session package for [fasthttp](https://github.com/valyala/fasthttp) servers

# Description

fasthttpsession is a session manager for Go. It only supports [fasthttp](https://github.com/valyala/fasthttp), currently support providers:

- file
- memcache
- memory
- mysql
- postgres
- redis
- sqlite3

# Features

- Focus on the design of the code architecture and expansion.
- Provide full session storage.
- Convenient switching of session storage.
- Customizable data serialization.
- Implement concurrent map(ccmap.go) to improve performance.

# Install

The only requirement is the Go Programming Language, at least v1.7

```shell
$ go get -u github.com/phachon/fasthttpsession
$ go get ./...
```

# Used

## Quick Start
```Golang

// fasthttpsession use memory provider

import (
	"github.com/phachon/fasthttpsession"
	"github.com/phachon/fasthttpsession/memory"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

// default config
var session = fasthttpsession.NewSession(fasthttpsession.NewDefaultConfig())

func main()  {
	// you must set up provider before use
	err := session.SetProvider("memory", &memory.Config{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	addr := ":8086"
	log.Println("fasthttpsession example server listen: "+addr)
	
	// fasthttp start listen serve
	err = fasthttp.ListenAndServe(addr, requestHandle)
	if err != nil {
		log.Println("listen server error :"+err.Error())
	}
}

// request handler
func requestHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("name", "fasthttpsession")

	ctx.SetBodyString(fmt.Sprintf("fasthttpsession setted key name= %s ok", sessionStore.Get("name").(string)))
}
```

## Custom configuration

If you don't want to use the default configuration, please use the following struct custom.
```Golang
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
```

Different session provider config, please look at the Config.go the provider name directory.

# Documents

Document address: [http://godoc.org/github.com/phachon/fasthttpsession](http://godoc.org/github.com/phachon/fasthttpsession)

# Example

[Some Example](_examples)

## Feedback

- If you like the project, please [Start](https://github.com/phachon/fasthttpsession/stargazers).
- If you have any problems in the process of use, welcome submit [Issue](https://github.com/phachon/fasthttpsession/issues).
- If you find and solve bug, welcome submit [Pull Request](https://github.com/phachon/fasthttpsession/pulls).
- If you want to expand session provider, welcome [Fork](https://github.com/phachon/fasthttpsession/network/members) and merge this rep.
- If you want to make a friend, welcome send email to [phachon@163.com](mailto:phachon@163.com).

## License

MIT

Thanks
---------
Create By phachon@163.com