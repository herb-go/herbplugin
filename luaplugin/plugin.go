package luaplugin

import (
	"path/filepath"
	"sync"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

const PluginType = "lua"
const DefaultNamespace = "system"

func New() *Plugin {
	return &Plugin{
		BasicPlugin: herbplugin.New(),
	}
}

type Plugin struct {
	sync.RWMutex
	entry  string
	LState *lua.LState
	*herbplugin.BasicPlugin
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
	p.BasicPlugin.MustInitPlugin()
	p.Builtin = map[string]lua.LGFunction{}
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].InitProcess != nil {
			processs = append(processs, p.modules[k].InitProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	if !p.DisableBuiltin {
		p.LState.PreloadModule(p.namespace, p.builtinLoader)
	}
}
func (p *Plugin) MustLoadPlugin(opt *herbplugin.Options) {
	p.BasicPlugin.MustLoadPlugin(opt)
	if p.entry != "" {
		err := p.LState.DoFile(filepath.Join(p.GetPluginLocation().Path, p.entry))
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plugin) builtinLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), p.Builtin)
	L.Push(mod)
	return 1
}
func (p *Plugin) MustBootPlugin() {
	p.BasicPlugin.MustBootPlugin()
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
	for k := range p.modules {
		if p.modules[k].CloseProcess != nil {
			processs = append(processs, p.modules[k].CloseProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	p.BasicPlugin.MustClosePlugin()
}
func (p *Plugin) LoadLuaPlugin() *Plugin {
	return p
}

type LuaPluginLoader interface {
	LoadLuaPlugin() *Plugin
}
