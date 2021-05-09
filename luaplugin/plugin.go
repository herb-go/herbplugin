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

func (p *Plugin) InitPlugin(opt *herbplugin.Options) error {
	err := p.BasicPlugin.InitPlugin(opt)
	if err != nil {
		return err
	}
	for _, v := range p.modules {
		err = v.InstallFn(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Plugin) StartPlugin() error {
	err := p.BasicPlugin.StartPlugin()
	if err != nil {
		return err
	}
	if p.startCommand != "" {
		err = p.LState.DoString(p.startCommand)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Plugin) LoadPlugin() error {
	if p.entry != "" {
		err := p.LState.DoFile(filepath.Join(p.GetPluginLocation().Path, p.entry))
		if err != nil {
			return err
		}
	}
	return p.BasicPlugin.LoadPlugin()
}
func (p *Plugin) ClosePlugin() error {
	for _, v := range p.modules {
		err := v.UninstallFn(p)
		if err != nil {
			return err
		}
	}
	p.LState.Close()
	return p.BasicPlugin.ClosePlugin()
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
