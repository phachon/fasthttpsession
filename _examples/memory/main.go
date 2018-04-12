package main

// fasthttpsession memory provider example

import (
	"github.com/phachon/fasthttpsession"
	"github.com/phachon/fasthttpsession/memory"
	"github.com/valyala/fasthttp"
	"log"
	"fmt"
	"os"
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

func main()  {

	// You must set up provider before use
	err := session.SetProvider("memory", &memory.Config{})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	addr := ":8086"
	log.Println("fasthttpsession memory example server listen: "+addr)
	// Fasthttp start listen serve
	err = fasthttp.ListenAndServe(addr, requestRouter)
	if err != nil {
		log.Println("listen server error :"+err.Error())
	}
}

// request router
func requestRouter(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		indexHandler(ctx)
	case "/set":
		setHandler(ctx)
	case "/get":
		getHandler(ctx)
	case "/delete":
		deleteHandle(ctx)
	case "/getAll":
		getAllHandle(ctx)
	case "/flush":
		flushHandle(ctx)
	case "/destroy":
		destroyHandle(ctx)
	case "/sessionid":
		sessionIdHandle(ctx)
	case "/regenerate":
		regenerateHandle(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

// index handler
func indexHandler(ctx *fasthttp.RequestCtx) {

	html := "<h2>Welcome to use fasthttpsession "+fasthttpsession.Version()+", you should request to the: </h2>"

	html += `> <a href="/">/</a><br>`
	html += `> <a href="/set">set</a><br>`
	html += `> <a href="/get">get</a><br>`
	html += `> <a href="/delete">delete</a><br>`
	html += `> <a href="/getAll">getAll</a><br>`
	html += `> <a href="/flush">flush</a><br>`
	html += `> <a href="/destroy">destroy</a><br>`
	html += `> <a href="/sessionid">sessionid</a><br>`
	html += `> <a href="/regenerate">regenerate</a><br>`

	ctx.SetContentType("text/html;charset=utf-8")
	ctx.SetBodyString(html)
}

// set handler
func setHandler(ctx *fasthttp.RequestCtx) {
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

// get handler
func getHandler(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	s := sessionStore.Get("name")
	if s == nil {
		ctx.SetBodyString("fasthttpsession get name is nil")
		return
	}

	ctx.SetBodyString(fmt.Sprintf("fasthttpsession get name= %s ok", s.(string)))
}

// delete handler
func deleteHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Delete("name")

	s := sessionStore.Get("name")
	if s == nil {
		ctx.SetBodyString("fasthttpsession delete key name ok")
		return
	}
	ctx.SetBodyString("fasthttpsession delete key name error")
}

// get all handler
func getAllHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("foo1", "baa1")
	sessionStore.Set("foo2", "baa2")
	sessionStore.Set("foo3", "baa3")
	sessionStore.Set("foo4", "baa5")

	data := sessionStore.GetAll()

	fmt.Println(data)
	ctx.SetBodyString("fasthttpsession get all data")
}

// flush handle
func flushHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Flush()

	ctx.SetBodyString("fasthttpsession flush data")
}

// destroy handle
func destroyHandle(ctx *fasthttp.RequestCtx) {
	// destroy session
	session.Destroy(ctx)

	ctx.SetBodyString("fasthttpsession destroy")
}

// get sessionId handle
func sessionIdHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionId := sessionStore.GetSessionId()
	ctx.SetBodyString("fasthttpsession sessionId: "+sessionId)
}

// regenerate handler
func regenerateHandle(ctx *fasthttp.RequestCtx) {
	// start session
	sessionStore, err := session.Start(ctx)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	// must defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("name", "foo")
	sessionStore.Get("name")

	sessionId := sessionStore.GetSessionId()

	ctx.SetBodyString("fasthttpsession regenerate sessionId: "+sessionId)
}
