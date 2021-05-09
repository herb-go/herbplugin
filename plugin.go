package herbplugin

type Plugin interface {
	InitPlugin(opt *Options) error
	LoadPlugin() error
	ClosePlugin() error
	StartPlugin() error
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

func (p *BasicPlugin) ClosePlugin() error {
	return nil
}

func (p *BasicPlugin) LoadPlugin() error {
	return nil
}
func (p *BasicPlugin) StartPlugin() error {
	p.context.Start()
	return nil
}
func (p *BasicPlugin) InitPlugin(opt *Options) error {
	p.context = NewContext()
	p.options = opt
	return nil
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
