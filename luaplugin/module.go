package luaplugin

import (
	"context"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

const ModuleIDParam = "builtin.param"
const ModuleIDPrint = "builtin.print"

func AppendCommonModules(i *Initializer) {
	i.Modules = append(i.Modules,
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
			L.Push(lua.LString(plugin.GetPluginParam(name)))
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
			plugin.PluginDebug(info)
			return 0
		}))
		next(ctx, plugin)
	},
	nil,
	nil,
)
