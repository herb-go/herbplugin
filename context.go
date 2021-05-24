package herbplugin

import "log"

func LogError(err error) {
	log.Println(err)
}

func NopDebuger(info string) {

}

type Context struct {
	components   map[string]interface{}
	errorHandler func(err error)
	debuger      func(info string)
}

func (c *Context) RegisterComponent(name string, component interface{}) {
	c.components[name] = component
}

func (c *Context) GetComponent(name string) (component interface{}, ok bool) {

	component, ok = c.components[name]
	return
}

func (c *Context) SetErrorHandler(h func(err error)) {
	c.errorHandler = h
}
func (c *Context) HandleError(err error) {
	c.errorHandler(err)
}

func (c *Context) SetDebuger(h func(info string)) {
	c.debuger = h
}
func (c *Context) Debug(info string) {
	c.debuger(info)
}

func NewContext() *Context {
	return &Context{
		components:   map[string]interface{}{},
		errorHandler: LogError,
		debuger:      NopDebuger,
	}
}
