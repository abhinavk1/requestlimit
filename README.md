# Gin Request Size Limit Middleware

Middleware for limiting size of POST requests for the 
[Gin-gonic Web Framework](https://github.com/gin-gonic/gin).

## Installation
This project uses [Go modules](https://blog.golang.org/using-go-modules).

```
go get -u github.com/abhinavk1/requestlimit
```

## Usage

### Quick Start
By default, the HTTP error code 413 will be sent back to the client if the request limit
is breached.
 
```go
package main

import (
	"github.com/abhinavk1/requestlimit"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handler(ctx *gin.Context) {
	val := ctx.PostForm("b")
	if len(ctx.Errors) > 0 {
		return
	}
	ctx.String(http.StatusOK, "got %s\n", val)
}

func main() {
	rtr := gin.Default()
	rtr.Use(requestlimit.Handler(10, nil)) // 10 bytes request size limit
	rtr.POST("/", handler)
	rtr.Run(":8080")
}
```

### Using a callback
You can also provided a callback method that allows you to choose what kind of HTTP 
response code and message to sent to the client.

The callback method needs to follow this signature,
```
func(ctx *gin.Context, err error)
```

For example,
```go
package main

import (
	"github.com/abhinavk1/requestlimit"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handler(ctx *gin.Context) {
	val := ctx.PostForm("b")
	if len(ctx.Errors) > 0 {
		return
	}
	ctx.String(http.StatusOK, "got %s\n", val)
}

func main() {
	rtr := gin.Default()
	rtr.Use(requestlimit.Handler(10, onRequestLimitReached)) // 10 bytes request size limit
	
	rtr.POST("/", handler)
	rtr.Run(":8081")
}

func onRequestLimitReached(ctx *gin.Context, err error) {
    ctx.Error(err)
    ctx.Header("connection", "close")
    ctx.String(http.StatusBadRequest, "request size was out of bounds")
    ctx.Abort()
}
```
