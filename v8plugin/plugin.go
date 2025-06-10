package v8plugin

import (
	"os"
	"path/filepath"
	"sync"

	v8 "rogchap.com/v8go"

	"github.com/herb-go/herbplugin"
)

const PluginType = "js"
const DefaultNamespace = "system"

func New() *Plugin {
	return &Plugin{
		Plugin: herbplugin.New(),
	}
}

func MustGetArg(call *v8.FunctionCallbackInfo, idx int) *v8.Value {
	args := call.Args()
	if idx < 0 || idx >= len(args) {
		return v8.Null(call.Context().Isolate())
	}
	return args[idx]
}
func MustNewValue(ctx *v8.Context, value interface{}) *v8.Value {
	switch value.(type) {
	case int:
		value = value.(int64)
	}
	v, err := v8.NewValue(ctx.Isolate(), value)
	if err != nil {
		panic(err)
	}
	return v
}
func MustObjectTemplateToValue(obj *v8.ObjectTemplate, ctx *v8.Context) *v8.Value {
	if obj == nil {
		return v8.Null(ctx.Isolate())
	}
	value, err := obj.NewInstance(ctx)
	if err != nil {
		panic(err)
	}
	return value.Value
}
func MustNewArray(ctx *v8.Context, args []v8.Valuer) *v8.Value {
	array, err := ctx.Global().Get("Array")
	if err != nil {
		panic(err)
	}
	fn, err := array.AsFunction()
	if err != nil {
		panic(err)
	}
	result, err := fn.Call(ctx.Global(), args...)
	if err != nil {
		panic(err)
	}
	return result
}
func MustSetObjectMethod(ctx *v8.Context, obj *v8.ObjectTemplate, name string, fn v8.FunctionCallback) {
	if obj == nil {
		return
	}
	method := v8.NewFunctionTemplate(ctx.Isolate(), fn)
	if method == nil {
		panic("Failed to create function template")
	}
	obj.Set(name, method)
}

type Plugin struct {
	sync.RWMutex
	entry   string
	Runtime *v8.Context
	herbplugin.Plugin
	DisableBuiltin bool
	startCommand   string
	modules        []*herbplugin.Module
	namespace      string
	Builtin        map[string]v8.FunctionCallback
}

func (p *Plugin) PluginType() string {
	return PluginType
}
func (p *Plugin) MustInitPlugin() {
	p.Plugin.MustInitPlugin()
	p.Builtin = map[string]v8.FunctionCallback{}
	var processs = make([]herbplugin.Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].InitProcess != nil {
			processs = append(processs, p.modules[k].InitProcess)
		}
	}
	builtin, err := v8.NewObjectTemplate(p.Runtime.Isolate()).NewInstance(p.Runtime)
	if err != nil {
		panic(err)
	}
	for key, fn := range p.Builtin {
		builtin.Set(key, fn)
	}
	herbplugin.Exec(p, processs...)
	if !p.DisableBuiltin {
		err := p.Runtime.Global().Set(p.namespace, builtin)
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
		_, err = p.Runtime.RunScript(string(data), p.entry)
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
		_, err := p.Runtime.RunScript(p.startCommand, "")
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
