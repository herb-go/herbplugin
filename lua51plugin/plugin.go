package lua51plugin

import (
	"path/filepath"
	"sync"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

const PluginType = "lua51"
const DefaultNamespace = "system"

func New() *Plugin {
	return &Plugin{
		Plugin: herbplugin.New(),
	}
}

type Plugin struct {
	sync.RWMutex
	entry  string
	LState *lua.LState
	herbplugin.Plugin
	DisableBuiltin bool
	startCommand   string
	modules        []*herbplugin.Module
	namespace      string
	Builtin        map[string]lua.LGFunction
}

func (p *Plugin) PluginType() string {
	return PluginType
}
func (p *Plugin) MustInitPlugin() {
	p.Plugin.MustInitPlugin()
	p.Builtin = map[string]lua.LGFunction{}
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].InitProcess != nil {
			processs = append(processs, p.modules[k].InitProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	if !p.DisableBuiltin {
		t := p.LState.NewTable()
		ft := p.LState.SetFuncs(t, p.Builtin)
		p.LState.SetGlobal(p.namespace, ft)
		// p.LState.PreloadModule(p.namespace, p.builtinLoader)
	}
}
func (p *Plugin) MustLoadPlugin(opt *herbplugin.Options) {
	p.Plugin.MustLoadPlugin(opt)
	if p.entry != "" {
		err := p.LState.DoFile(filepath.Join(p.GetPluginLocation().Path, p.entry))
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plugin) MustBootPlugin() {
	p.Plugin.MustBootPlugin()
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].BootProcess != nil {
			processs = append(processs, p.modules[k].BootProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	if p.startCommand != "" {
		err := p.LState.DoString(p.startCommand)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plugin) MustClosePlugin() {
	defer p.LState.Close()
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for i := len(p.modules) - 1; i >= 0; i-- {
		if p.modules[i].CloseProcess != nil {
			processs = append(processs, p.modules[i].CloseProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	p.Plugin.MustClosePlugin()
}
func (p *Plugin) LoadLuaPlugin() *Plugin {
	return p
}

type LuaPluginLoader interface {
	LoadLuaPlugin() *Plugin
}
