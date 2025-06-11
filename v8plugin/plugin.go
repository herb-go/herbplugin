package v8plugin

import (
	"errors"
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
	switch v := value.(type) {
	case *v8.Object:
		return v.Value
	case *v8.Value:
		return v
	case int:
		value = int64(v)
	}
	val, err := v8.NewValue(ctx.Isolate(), value)
	if err != nil {
		panic(err)
	}
	return val
}
func MustObjectTemplateToValue(ctx *v8.Context, obj *v8.ObjectTemplate) *v8.Value {
	if obj == nil {
		return v8.Null(ctx.Isolate())
	}
	value, err := obj.NewInstance(ctx)
	if err != nil {
		panic(err)
	}
	return value.Value
}
func MustNewArray(ctx *v8.Context, args []*v8.Value) *v8.Value {
	array, err := ctx.Global().Get("Array")
	if err != nil {
		panic(err)
	}
	fn, err := array.AsFunction()
	if err != nil {
		panic(err)
	}
	fnargs := make([]v8.Valuer, len(args))
	for i, v := range args {
		fnargs[i] = v
	}
	result, err := fn.Call(array, fnargs...)
	if err != nil {
		panic(err)
	}
	return result
}
func MustConvertToArray(ctx *v8.Context, val *v8.Value) []*v8.Value {
	if val.IsNull() || val.IsUndefined() {
		return []*v8.Value{}
	}
	if !val.IsArray() {
		panic(errors.New("value is not an array"))
	}

	obj, err := val.AsObject()
	if err != nil {
		panic(err)
	}
	length, err := obj.Get("length")
	defer length.Release()
	if err != nil {
		panic(err)
	}
	l := length.Int32()
	if l < 0 {
		panic(errors.New("array length is negative"))
	}
	result := make([]*v8.Value, l)
	for i := int32(0); i < l; i++ {
		result[int(i)], err = obj.GetIdx(uint32(i))
		if err != nil {
			panic(err)
		}
	}
	return result
}
func MustConvertToStringArray(ctx *v8.Context, val *v8.Value) []string {
	values := MustConvertToArray(ctx, val)
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
func WrapCallback(fn v8.FunctionCallback) v8.FunctionCallback {
	return func(info *v8.FunctionCallbackInfo) *v8.Value {
		return fn(info)
	}
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
func MustGetItem(obj *v8.Object, name string) *v8.Value {
	value, err := obj.Get(name)
	if err != nil {
		panic(err)
	}
	return value
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
	p.Runtime.Close()
	p.Runtime.Isolate().Dispose()
}
func (p *Plugin) LoadJsPlugin() *Plugin {
	return p
}

type JsPluginLoader interface {
	LoadJsPlugin() *Plugin
}
