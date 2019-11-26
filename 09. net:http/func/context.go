package net

import (
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"path/filepath"
)


// Params 필드에 라우터에서 해석한 URL 매개변수로 담고, 핸들러 내부에는 Context 값이 전달되게 한다.
type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request *http.Request
}

type HandlerFunc func(*Context)

func (c *Context) RenderJson(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(c.ResponseWriter).Encode(v); err!=nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderXml(v interface{}) {
	c.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	c.ResponseWriter.Header().Set("Content-Type", "application/xml; charset=utf-8")

	if err := xml.NewEncoder(c.ResponseWriter).Encode(v); err!=nil {
		c.RenderErr(http.StatusInternalServerError, err)
	}
}

func (c *Context) RenderErr(code int, err error) {
	if err!=nil {
		if code>0 {
			http.Error(c.ResponseWriter, http.StatusText(code), code)
		} else {
			defaultErr := http.StatusInternalServerError
			http.Error(c.ResponseWriter, http.StatusText(defaultErr), defaultErr)
		}
	}
}

var templates = map[string]*template.Template{}

func (c *Context) RenderTemplate(path string, v interface{}) {
	t, ok := templates[path]
	if !ok {
		t = template.Must(template.ParseFiles(filepath.Join(".", path)))
		templates[path] = t
	}

	t.Execute(c.ResponseWriter, v)
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, http.StatusMovedPermanently)
}