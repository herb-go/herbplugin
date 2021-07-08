package jsplugin

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/dop251/goja"
	"github.com/robertkrimen/otto"

	"github.com/herb-go/herbplugin"
)

const PluginType = "js"
const DefaultNamespace = "system"

func New() *Plugin {
	return &Plugin{
		Plugin: herbplugin.New(),
	}
}

type Plugin struct {
	sync.RWMutex
	entry   string
	Runtime *goja.Runtime
	herbplugin.Plugin
	DisableBuiltin bool
	startCommand   string
	modules        []*herbplugin.Module
	namespace      string
	Builtin        map[string]func(call otto.FunctionCall) otto.Value
}

func (p *Plugin) PluginType() string {
	return PluginType
}
func (p *Plugin) MustInitPlugin() {
	p.Plugin.MustInitPlugin()
	p.Builtin = map[string]func(call otto.FunctionCall) otto.Value{}
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].InitProcess != nil {
			processs = append(processs, p.modules[k].InitProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	if !p.DisableBuiltin {
		err := p.Runtime.Set(p.namespace, p.Builtin)
		if err != nil {
			panic(err)
		}
	}
}
func (p *Plugin) MustLoadPlugin() {
	p.Plugin.MustLoadPlugin()
	if p.entry != "" {
		data, err := os.ReadFile(filepath.Join(p.PluginOptions().GetLocation().Path, p.entry))
		if err != nil {
			panic(err)
		}
		_, err = p.Runtime.RunScript(p.entry, string(data))
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
		_, err := p.Runtime.RunString(p.startCommand)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plugin) MustClosePlugin() {
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for i := len(p.modules) - 1; i >= 0; i-- {
		if p.modules[i].CloseProcess != nil {
			processs = append(processs, p.modules[i].CloseProcess)
		}
	}
	herbplugin.Exec(p, processs...)
	p.modules = nil
	p.Builtin = nil
	p.Plugin.MustClosePlugin()
	p.Runtime = nil
}
func (p *Plugin) LoadJsPlugin() *Plugin {
	return p
}

type JsPluginLoader interface {
	LoadJsPlugin() *Plugin
}
