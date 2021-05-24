package herbplugin

type Plugin interface {
	MustInitPlugin()
	MustLoadPlugin(opt *Options)
	MustBootPlugin()
	MustClosePlugin()
	PluginType() string
	GetPluginParam(name string) string
	GetPluginTrusted() *Trusted
	GetPluginLocation() *Location
	GetPluginComponent(name string) (interface{}, bool)
	RegisterPluginComonent(name string, comonent interface{})
	SetPluginErrorHandler(func(err error))
	HandlePluginError(err error)
	SetPluginDebuger(func(info string))
	PluginDebug(info string)
}

type BasicPlugin struct {
	options *Options
	context *Context
}

func (p *BasicPlugin) MustInitPlugin() {
	p.context = NewContext()
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
func (p *BasicPlugin) GetPluginComponent(name string) (interface{}, bool) {
	return p.context.GetComponent(name)
}
func (p *BasicPlugin) RegisterPluginComonent(name string, component interface{}) {
	p.context.RegisterComponent(name, component)
}

func (p *BasicPlugin) SetPluginErrorHandler(h func(err error)) {
	p.context.SetErrorHandler(h)
}
func (p *BasicPlugin) HandlePluginError(err error) {
	p.context.HandleError(err)
}
func (p *BasicPlugin) SetPluginDebuger(h func(info string)) {
	p.context.SetDebuger(h)
}
func (p *BasicPlugin) PluginDebug(info string) {
	p.context.Debug(info)
}

func Lanuch(p Plugin, opt *Options) {
	p.MustInitPlugin()
	p.MustLoadPlugin(opt)
	p.MustBootPlugin()
}
