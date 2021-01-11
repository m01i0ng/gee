package gee

import (
  "log"
  "net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
  *RouterGroup
  router *router
  groups []*RouterGroup
}

func New() *Engine {
  e := &Engine{router: newRouter()}
  e.RouterGroup = &RouterGroup{engine: e}
  e.groups = []*RouterGroup{e.RouterGroup}
  return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  c := newContext(w, req)
  e.router.handle(c)
}

func (e *Engine) Run(addr string) error {
  return http.ListenAndServe(addr, e)
}

type RouterGroup struct {
  prefix      string
  middlewares []HandlerFunc
  parent      *RouterGroup
  engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
  e := group.engine
  newGroup := &RouterGroup{
    prefix: group.prefix + prefix,
    parent: group,
    engine: e,
  }
  e.groups = append(e.groups, newGroup)
  return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
  pattern := group.prefix + comp
  log.Printf("Route %4s -%s", method, pattern)
  group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
  group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
  group.addRoute("POST", pattern, handler)
}
