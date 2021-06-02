package lua51plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

const ModuleIDParam = "builtin.param"
const ModuleIDPrint = "builtin.print"
const ModuleIDOpenlib = "builtin.openlib"

func AppendCommonModules(i *Initializer) {
	i.Modules = append(i.Modules,
		ModuleOpenlib,
		ModuleParam,
		ModulePrint,
	)
}

var ModuleParam = herbplugin.CreateModule(
	ModuleIDParam,
	func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
		plugin := p.(LuaPluginLoader).LoadLuaPlugin()
		plugin.Builtin["getparam"] = func(L *lua.LState) int {
			name := L.ToString(1)
			L.Push(lua.LString(p.PluginOptions().GetParam(name)))
			return 1
		}
		next(ctx, plugin)
	},
	nil,
	nil,
)
var ModulePrint = herbplugin.CreateModule(
	ModuleIDPrint,
	func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
		plugin := p.(LuaPluginLoader).LoadLuaPlugin()
		plugin.LState.SetGlobal("print", plugin.LState.NewFunction(func(L *lua.LState) int {
			info := L.ToString(1)
			plugin.PluginPrint(info)
			return 0
		}))
		next(ctx, plugin)
	},
	nil,
	nil,
)

var safetycommands = []string{
	"os.remove=nil",
	"os.rename=nil",
	"os.execute=nil",
	"os.getenv=nil",
	"os.setenv=nil",
	"os.tmpname=nil",
	"os.exit=nil",
	"os.setlocale=nil",
	"io=nil",
	"dofile=nil",
	"loadfile=nil",
	"load=nil",
}

func pluginDoFile(p *Plugin) func(L *lua.LState) int {
	location := p.PluginOptions().GetLocation()
	return func(L *lua.LState) int {
		src := L.ToString(1)
		top := L.GetTop()
		cleanpath := location.MustCleanInsidePath(src)
		if cleanpath == "" {
			L.RaiseError("%s not in script location", src)
		}
		fn, err := L.LoadFile(cleanpath)
		if err != nil {
			L.Push(lua.LString(err.Error()))
			L.Panic(L)
		}
		L.Push(fn)
		L.Call(0, lua.MultRet)
		return L.GetTop() - top
	}
}
func pluginLoadFile(p *Plugin) func(L *lua.LState) int {
	location := p.PluginOptions().GetLocation()
	return func(L *lua.LState) int {
		src := L.ToString(1)
		cleanpath := location.MustCleanInsidePath(src)
		if cleanpath == "" {
			L.RaiseError("%s not in script location", src)
		}
		fn, err := L.LoadFile(cleanpath)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(fn)
		return 1
	}
}
func pluginLoaders(p *Plugin) []func(L *lua.LState) int {
	location := p.PluginOptions().GetLocation()
	return []func(L *lua.LState) int{
		loLoaderPreload,
		func(L *lua.LState) int {
			name := L.CheckString(1)
			cleanpath := strings.Replace(name, ".", string(os.PathSeparator), -1)
			cleanpath += ".lua"
			cleanpath = location.MustCleanInsidePath(cleanpath)

			if cleanpath == "" {
				L.RaiseError("%s not in script location", name)
			}
			fn, err1 := L.LoadFile(cleanpath)
			if err1 != nil {
				L.RaiseError(err1.Error())
			}
			L.Push(fn)
			return 1
		},
	}
}

func loLoaderPreload(L *lua.LState) int {
	name := L.CheckString(1)
	preload := L.GetField(L.GetField(L.Get(lua.EnvironIndex), "package"), "preload")
	if _, ok := preload.(*lua.LTable); !ok {
		L.RaiseError("package.preload must be a table")
	}
	lv := L.GetField(preload, name)
	if lv == lua.LNil {
		L.Push(lua.LString(fmt.Sprintf("no field package.preload['%s']", name)))
		return 1
	}
	L.Push(lv)
	return 1
}

var ModuleOpenlib = herbplugin.CreateModule(
	ModuleIDOpenlib,
	func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
		plugin := p.(LuaPluginLoader).LoadLuaPlugin()
		plugin.LState.OpenLibs()

		if !plugin.PluginOptions().MustAuthorizePermission(herbplugin.PermissionDangerousAPI) {
			for _, v := range safetycommands {
				err := plugin.LState.DoString(v)
				if err != nil {
					panic(err)
				}
			}
			loLoaders := pluginLoaders(plugin)
			loaders := plugin.LState.CreateTable(len(loLoaders), 0)
			for i, loader := range loLoaders {
				plugin.LState.RawSetInt(loaders, i+1, plugin.LState.NewFunction(loader))
			}
			plugin.LState.SetField(plugin.LState.Get(lua.RegistryIndex), "_LOADERS", loaders)
			plugin.LState.SetGlobal("dofile", plugin.LState.NewFunction(pluginDoFile(plugin)))
			plugin.LState.SetGlobal("loadfile", plugin.LState.NewFunction(pluginLoadFile(plugin)))
		}
		next(ctx, plugin)
	},
	nil,
	nil,
)
