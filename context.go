package herbplugin

import "log"

func LogError(err error) {
	log.Println(err)
}

func NopPrinter(info string) {
}

type Context struct {
	errorHandler func(err error)
	printer      func(info string)
}

func (c *Context) SetErrorHandler(h func(err error)) {
	c.errorHandler = h
}
func (c *Context) HandleError(err error) {
	c.errorHandler(err)
}

func (c *Context) SetPrinter(h func(info string)) {
	c.printer = h
}
func (c *Context) Print(info string) {
	c.printer(info)
}

func NewContext() *Context {
	return &Context{
		errorHandler: LogError,
		printer:      NopPrinter,
	}
}
