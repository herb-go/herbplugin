package luaplugin

import (
	"context"

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

var ModuleParam = CreateModule(
	ModuleIDParam,
	func(ctx context.Context, plugin *Plugin, next func(ctx context.Context, plugin *Plugin)) {
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
var ModulePrint = CreateModule(
	ModuleIDPrint,
	func(ctx context.Context, plugin *Plugin, next func(ctx context.Context, plugin *Plugin)) {
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

type Module struct {
	ID           string
	InitProcess  Process
	BootProcess  Process
	CloseProcess Process
}

func CreateModule(id string, initfn Process, bootfn Process, closefn Process) *Module {
	return &Module{
		ID:           id,
		InitProcess:  initfn,
		BootProcess:  bootfn,
		CloseProcess: closefn,
	}
}
