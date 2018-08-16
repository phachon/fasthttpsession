[![logo](./logo.png)](https://github.com/savsgio/fasthttpsession)

[![build](https://img.shields.io/shippable/5444c5ecb904a4b21567b0ff.svg)](https://travis-ci.org/phachon/fasthttpsession)
[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/savsgio/fasthttpsession)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phachon/fasthttpsession/master/LICENSE)
[![go_Report](https://goreportcard.com/badge/github.com/savsgio/fasthttpsession)](https://goreportcard.com/report/github.com/savsgio/fasthttpsession)
[![release](https://img.shields.io/github/release/phachon/fasthttpsession.svg?style=flat)](https://github.com/savsgio/fasthttpsession/releases) 
[![powered_by](https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat)]()
[![platforms](https://img.shields.io/badge/platform-All-yellow.svg?style=flat)]()

fasthttpsession 是一个快速且强大的 [fasthttp](https://github.com/valyala/fasthttp) session 管理包

[English Document](./README.md)

# 描述

fasthttpsession 是 Go 实现的一个 session 管理器。它只能用于 [fasthttp](https://github.com/valyala/fasthttp) 框架, 目前支持的 session 存储如下:

- file
- memcache
- memory
- mysql
- postgres
- redis
- sqlite3

# 功能

- 关注代码架构和扩展的设计。
- 提供全面的 session 存储。
- 方便的 session 存储切换。
- 可自由定制数据序列化函数。
- 实现了并发 map(ccmap.go) 去提高性能。

# 安装

要求是 Go 至少是 v1.7。

```shell
$ go get -u github.com/savsgio/fasthttpsession
$ go get ./...
```

# 使用

## 快速开始
```Golang

// fasthttpsession use memory provider

import (
	"github.com/savsgio/fasthttpsession"
	"github.com/savsgio/fasthttpsession/memory"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

// 默认的 session 全局配置
var session = fasthttpsession.NewSession(fasthttpsession.NewDefaultConfig())

func main()  {
	// 必须在使用之前指定 session 的存储
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
	// 必须 defer sessionStore.save(ctx)
	defer sessionStore.Save(ctx)

	sessionStore.Set("name", "fasthttpsession")

	ctx.SetBodyString(fmt.Sprintf("fasthttpsession setted key name= %s ok", sessionStore.Get("name").(string)))
}
```

## 自定义配置

如果您不想使用默认配置，请使用以下结构自定义你想要的配置。
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
	
	// sessionID is in url query
	SessionIdInURLQuery bool
	
	// sessionName in url query
	SessionNameInUrlQuery string
	
	// sessionID is in http header
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

不同的 session 存储提供有着不同的配置，请查看存储名称目录下的 Config.go

# 文档

文档地址: [http://godoc.org/github.com/savsgio/fasthttpsession](http://godoc.org/github.com/savsgio/fasthttpsession)

# 示例

[一些例子](_examples)

## 反馈

- 如果您喜欢该项目，请 [Start](https://github.com/savsgio/fasthttpsession/stargazers).
- 如果在使用过程中有任何问题， 请提交 [Issue](https://github.com/savsgio/fasthttpsession/issues).
- 如果您发现并解决了bug，请提交 [Pull Request](https://github.com/savsgio/fasthttpsession/pulls).
- 如果您想扩展 session 存储，欢迎 [Fork](https://github.com/savsgio/fasthttpsession/network/members) and merge this rep.
- 如果你想交个朋友，欢迎发邮件给 [phachon@163.com](mailto:phachon@163.com).

## License

MIT

Thanks
---------
Create By phachon@163.com
