package main

import (
  "log"
  "net/http"
  "time"

  "github.com/m01i0ng/gee"
)

func onlyForV2() gee.HandlerFunc {
  return func(c *gee.Context) {
    t := time.Now()
    c.Fail(500, "Internal Server Error")
    log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
  }
}

func main() {
  g := gee.New()
  g.Use(gee.Logger())
  g.GET("/", func(c *gee.Context) {
    c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
  })
  v2 := g.Group("/v2")
  v2.Use(onlyForV2())
  {
    v2.GET("/hello/:name", func(c *gee.Context) {
      c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
    })
  }
  g.Run(":9999")
}
