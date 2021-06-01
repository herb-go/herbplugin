package herbplugin

import "log"

type Plugin interface {
	MustInitPlugin()
	MustLoadPlugin()
	MustBootPlugin()
	MustClosePlugin()
	PluginType() string
	SetPluginErrorHandler(func(err error))
	HandlePluginError(err error)
	SetPluginPrinter(func(info string))
	PluginPrint(info string)
	PluginOptions() Options
	SetPluginOptions(opt Options)
}

func LogError(err error) {
	log.Println(err)
}

func NopPrinter(info string) {
}

type BasicPlugin struct {
	options      Options
	errorHandler func(err error)
	printer      func(info string)
}

func (p *BasicPlugin) MustInitPlugin() {
}
func (p *BasicPlugin) MustClosePlugin() {

}

func (p *BasicPlugin) MustLoadPlugin() {
}

func (p *BasicPlugin) MustBootPlugin() {

}
func (p *BasicPlugin) PluginType() string {
	return "unknown"
}

func (p *BasicPlugin) SetPluginErrorHandler(h func(err error)) {
	p.errorHandler = h
}
func (p *BasicPlugin) HandlePluginError(err error) {
	p.errorHandler(err)
}
func (p *BasicPlugin) SetPluginPrinter(h func(info string)) {
	p.printer = (h)
}
func (p *BasicPlugin) PluginPrint(info string) {
	p.printer(info)
}
func (p *BasicPlugin) PluginOptions() Options {
	return p.options
}
func (p *BasicPlugin) SetPluginOptions(opt Options) {
	p.options = opt
}
func New() *BasicPlugin {
	p := &BasicPlugin{}
	p.errorHandler = LogError
	p.printer = NopPrinter
	return p
}

func Lanuch(p Plugin, opt Options) {
	p.SetPluginOptions(opt)
	p.MustInitPlugin()
	p.MustLoadPlugin()
	p.MustBootPlugin()
}
