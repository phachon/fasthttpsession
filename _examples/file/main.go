package main

// fasthttpsession file provider example

import (
	"log"
	"os"

	"github.com/savsgio/fasthttpsession"
	"github.com/savsgio/fasthttpsession/file"
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
	err := session.SetProvider("file", &file.Config{
		SavePath: ".session", // session file save path
	})

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	addr := ":8086"
	log.Println("fasthttpsession file example server listen: " + addr)
	// Fasthttp start listen serve
	err = fasthttp.ListenAndServe(addr, requestRouter)
	if err != nil {
		log.Println("listen server error :" + err.Error())
	}
}
