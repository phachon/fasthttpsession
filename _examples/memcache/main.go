package main

// fasthttpsession redis provider example

import (
	"log"
	"os"

	"github.com/phachon/fasthttpsession"
	"github.com/phachon/fasthttpsession/memcache"
	"github.com/valyala/fasthttp"
)

// default config
var session = fasthttpsession.NewSession(fasthttpsession.NewDefaultConfig())

// custom config
//var session = fasthttpsession.NewSession(&fasthttpsession.Config{
//	CookieName: "ssid",
//	Domain: "",
//	Expires: time.Hour * 2,
//	GCLifetime: 3,
//	SessionLifetime: 60,
//	Secure: true,
//	SessionIdInURLQuery: false,
//	SessionNameInUrlQuery: "",
//	SessionIdInHttpHeader: false,
//	SessionNameInHttpHeader: "",
//	SessionIdGeneratorFunc: func() string {return ""},
//	EncodeFunc: func(cookieValue string) (string, error) {return "", nil},
//	DecodeFunc: func(cookieValue string) (string, error) {return "", nil},
//})

func main() {

	// You must set up provider before use
	err := session.SetProvider("memcache", &memcache.Config{
		ServerList: []string{
			"127.0.0.1:21122",
			"127.0.0.1:21123",
		},
		MaxIdle:   8,
		KeyPrefix: "session",
	})

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	addr := ":8086"
	log.Println("fasthttpsession redis example server listen: " + addr)
	// Fasthttp start listen serve
	err = fasthttp.ListenAndServe(addr, requestRouter)
	if err != nil {
		log.Println("listen server error :" + err.Error())
	}
}
