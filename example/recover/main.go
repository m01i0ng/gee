package main

import (
  "net/http"

  "github.com/m01i0ng/gee"
)

func main() {
  g := gee.Default()
  g.GET("/", func(c *gee.Context) {
    c.String(http.StatusOK, "Hello M\n")
  })
  g.GET("/panic", func(c *gee.Context) {
    names := []string{"M"}
    c.String(http.StatusOK, names[100])
  })
  g.Run(":9999")
}
