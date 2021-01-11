package gee

import (
  "html/template"
  "log"
  "net/http"
  "path"
  "strings"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
  *RouterGroup
  router        *router
  groups        []*RouterGroup
  htmlTemplates *template.Template
  funcMap       template.FuncMap
}

func New() *Engine {
  e := &Engine{router: newRouter()}
  e.RouterGroup = &RouterGroup{engine: e}
  e.groups = []*RouterGroup{e.RouterGroup}
  return e
}

func Default() *Engine {
  e := New()
  e.Use(Logger(), Recovery())
  return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  var middlewares []HandlerFunc
  for _, group := range e.groups {
    if strings.HasPrefix(req.URL.Path, group.prefix) {
      middlewares = append(middlewares, group.middlewares...)
    }
  }
  c := newContext(w, req)
  c.handlers = middlewares
  c.engine = e
  e.router.handle(c)
}

func (e *Engine) Run(addr string) error {
  return http.ListenAndServe(addr, e)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
  e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
  e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).Parse(pattern))
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

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
  group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
  absolutePath := path.Join(group.prefix, relativePath)
  fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
  return func(c *Context) {
    file := c.Param("filepath")
    if _, err := fs.Open(file); err != nil {
      c.Status(http.StatusNotFound)
      return
    }
    fileServer.ServeHTTP(c.Writer, c.Req)
  }
}

func (group *RouterGroup) Static(relativePath string, root string) {
  handler := group.createStaticHandler(relativePath, http.Dir(root))
  urlPattern := path.Join(relativePath, "/*filepath")
  group.GET(urlPattern, handler)
}
