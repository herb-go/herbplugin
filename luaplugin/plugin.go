package luaplugin

import (
	"path/filepath"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

func New() *Plugin {
	return &Plugin{
		LState: lua.NewState(),
	}
}

type Plugin struct {
	entry  string
	LState *lua.LState
	herbplugin.BasicPlugin
	startCommand string
	modules      []*Module
}

func (p *Plugin) MustInitPlugin(opt *herbplugin.Options) error {
	p.BasicPlugin.MustInitPlugin(opt)
	for _, v := range p.modules {
		err := v.InstallFn(p)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (p *Plugin) MustStartPlugin() {
	p.BasicPlugin.MustStartPlugin()

	if p.startCommand != "" {
		err := p.LState.DoString(p.startCommand)
		if err != nil {
			panic(err)
		}
	}
}
func (p *Plugin) MustLoadPlugin() {
	if p.entry != "" {
		err := p.LState.DoFile(filepath.Join(p.GetPluginLocation().Path, p.entry))
		if err != nil {
			panic(err)
		}
	}
	p.BasicPlugin.MustLoadPlugin()
}
func (p *Plugin) MustClosePlugin() {
	for _, v := range p.modules {
		err := v.UninstallFn(p)
		if err != nil {
			panic(err)
		}
	}
	p.LState.Close()
	p.BasicPlugin.MustClosePlugin()
}

type Module struct {
	Name        string
	InstallFn   func(*Plugin) error
	UninstallFn func(*Plugin) error
}

func CreateModule(name string, installfn func(*Plugin) error, uninstallfn func(*Plugin) error) *Module {
	return &Module{
		Name:        name,
		InstallFn:   installfn,
		UninstallFn: uninstallfn,
	}
}

type PluginFactory struct {
	Entry             string
	StartCommand      string
	RegisteredModules []*Module
}

func (f *PluginFactory) CreatePlugin() *Plugin {
	p := New()
	p.entry = f.Entry
	p.modules = f.RegisteredModules
	p.startCommand = f.StartCommand
	return p
}
