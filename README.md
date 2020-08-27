# Gin Request Size Limit Middleware

Limit size of POST requests for Gin framework

## Example

[embedmd]:# (example/main.go go)
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
	// Default configuration
	rtr := gin.Default()
	rtr.Use(requestlimit.Handler(10, nil))
	rtr.POST("/", handler)
	rtr.Run(":8080")

	// With custom limit reached hook
	rtr2 := gin.Default()
	rtr2.Use(requestlimit.Handler(10, func(ctx *gin.Context, err error) {
		ctx.Error(err)
		ctx.Header("connection", "close")
		ctx.String(http.StatusRequestEntityTooLarge, "request size was out of bounds")
		ctx.Abort()
	}))
	
	rtr2.POST("/", handler)
	rtr2.Run(":8081")
}

```
