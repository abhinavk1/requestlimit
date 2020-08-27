package requestlimit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestSizeLimiterOK(t *testing.T) {
	router := gin.New()
	router.Use(Handler(10, nil))
	router.POST("/test_ok", func(c *gin.Context) {
		ioutil.ReadAll(c.Request.Body)
		if len(c.Errors) > 0 {
			return
		}
		c.Request.Body.Close()
		c.String(http.StatusOK, "OK")
	})
	resp := performRequest(http.MethodPost, "/test_ok", "big=abc", router)

	if resp.Code != http.StatusOK {
		t.Fatalf("error posting - http status %v", resp.Code)
	}
}

func TestRequestSizeLimiterOver(t *testing.T) {
	router := gin.New()
	router.Use(Handler(10, nil))
	router.POST("/test_large", testLimitOver)
	resp := performRequest(http.MethodPost, "/test_large", "big=abcdefghijklmnop", router)

	if resp.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("error posting - http status %v", resp.Code)
	}

	fmt.Println(resp.Body.String())

	t.Run("custom onLimitReached function", func(t *testing.T) {

		router = gin.New()
		router.Use(Handler(10, func(ctx *gin.Context, err error) {
			ctx.Error(err)
			ctx.Header("connection", "close")
			ctx.String(http.StatusBadRequest, "uploaded file was too large to process")
			ctx.Abort()
		}))
		router.POST("/test_large", testLimitOver)

		resp = performRequest(http.MethodPost, "/test_large", "big=abcdefghijklmnop", router)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("error posting - http status %v", resp.Code)
		}

		if resp.Body.String() != "uploaded file was too large to process" {
			t.Fatalf("error posting - http message %v", resp.Body.String())
		}
	})
}

func performRequest(method, target, body string, router *gin.Engine) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != "" {
		buf = new(bytes.Buffer)
		buf.WriteString(body)
	}
	r := httptest.NewRequest(method, target, buf)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

func testLimitOver(c *gin.Context) {
	ioutil.ReadAll(c.Request.Body)
	if len(c.Errors) > 0 {
		return
	}
	c.Request.Body.Close()
	c.String(http.StatusOK, "OK")
}
