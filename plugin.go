package herbplugin

import "log"

type Plugin interface {
	MustInitPlugin()
	MustLoadPlugin(opt *Options)
	MustBootPlugin()
	MustClosePlugin()
	PluginType() string
	GetPluginParam(name string) string
	GetPluginLocation() *Location
	SetPluginErrorHandler(func(err error))
	HandlePluginError(err error)
	SetPluginPrinter(func(info string))
	PluginPrint(info string)
	SetPluginAuthorizer(Authorizer)
	PluginAuthorizer() Authorizer
}

func LogError(err error) {
	log.Println(err)
}

func NopPrinter(info string) {
}

type BasicPlugin struct {
	options      *Options
	errorHandler func(err error)
	printer      func(info string)
	authorizer   Authorizer
}

func (p *BasicPlugin) MustInitPlugin() {
}
func (p *BasicPlugin) MustClosePlugin() {

}

func (p *BasicPlugin) MustLoadPlugin(opt *Options) {
	p.options = opt
}

func (p *BasicPlugin) MustBootPlugin() {

}
func (p *BasicPlugin) PluginType() string {
	return "unknown"
}
func (p *BasicPlugin) GetPluginParam(name string) string {
	return p.options.Params[name]
}
func (p *BasicPlugin) GetPluginTrusted() *Trusted {
	return p.options.Trusted
}
func (p *BasicPlugin) GetPluginLocation() *Location {
	return p.options.Location
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
func (p *BasicPlugin) SetPluginAuthorizer(a Authorizer) {
	p.authorizer = a
}
func (p *BasicPlugin) PluginAuthorizer() Authorizer {
	return p.authorizer
}

func New() *BasicPlugin {
	p := &BasicPlugin{}
	p.options = NewOptions()
	p.errorHandler = LogError
	p.printer = NopPrinter
	p.authorizer = NewBasicPluginAuthorizer(p)
	return p
}

func Lanuch(p Plugin, opt *Options) {
	p.MustInitPlugin()
	p.MustLoadPlugin(opt)
	p.MustBootPlugin()
}
