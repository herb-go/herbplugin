package herbplugin

type Plugin interface {
	MustInitPlugin(opt *Options)
	MustLoadPlugin()
	MustClosePlugin()
	MustStartPlugin()
	GetPluginParam(name string) string
	GetPluginTrusted() *Trusted
	GetPluginLocation() *Location
	GetPluginComponent(name string) (interface{}, bool)
	RegisterPluginComonent(name string, comonent interface{})
}

type BasicPlugin struct {
	options *Options
	context *Context
}

func (p *BasicPlugin) MustClosePlugin() {

}

func (p *BasicPlugin) MustLoadPlugin() {
}
func (p *BasicPlugin) MustStartPlugin() {
	p.context.Start()
}
func (p *BasicPlugin) MustInitPlugin(opt *Options) {
	p.context = NewContext()
	p.options = opt
}
func (p *BasicPlugin) GetPluginParam(name string) string {
	return p.options.Params[name]
}
func (p *BasicPlugin) GetPluginTrusted() *Trusted {
	return &p.options.Trusted
}
func (p *BasicPlugin) GetPluginLocation() *Location {
	return &p.options.Location
}
func (p *BasicPlugin) GetPluginComponent(name string) (interface{}, bool) {
	return p.context.GetComponent(name)
}
func (p *BasicPlugin) RegisterPluginComonent(name string, component interface{}) {
	p.context.RegisterComponent(name, component)
}
