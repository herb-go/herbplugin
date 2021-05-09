package herbplugin

import (
	"log"
)

func LogError(err error) {
	log.Println(err)
}

type Location struct {
	Path string
}

type Params map[string]string

type Trusted struct {
	Paths   []string
	Domains []string
}

type Options struct {
	Location
	Params
	Trusted
}

type Context struct {
	started      bool
	components   map[string]interface{}
	errorHandler func(err error)
}

func (c *Context) RegisterComponent(name string, component interface{}) {
	if c.started {
		panic(ErrPluginStarted)
	}
	c.components[name] = component
}

func (c *Context) GetComponent(name string) (component interface{}, ok bool) {
	if c.started {
		panic(ErrPluginNotStarted)
	}
	component, ok = c.components[name]
	return
}

func (c *Context) Start() {
	c.started = true
}

func NewContext() *Context {
	return &Context{
		components:   map[string]interface{}{},
		errorHandler: LogError,
	}
}
