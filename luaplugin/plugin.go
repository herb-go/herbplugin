package luaplugin

import (
	"context"
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

func (p *Plugin) MustInitPlugin(opt *herbplugin.Options) {
	p.BasicPlugin.MustInitPlugin(opt)
	var processs = make([]Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k] != nil {
			processs = append(processs, p.modules[k].InstallProcess)
		}
	}
	ComposeProcess(processs...)(context.TODO(), p, Nop)
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
	defer p.LState.Close()
	var processs = make([]Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k] != nil {
			processs = append(processs, p.modules[k].UninstallProcess)
		}
	}
	ComposeProcess(processs...)(context.TODO(), p, Nop)
	p.BasicPlugin.MustClosePlugin()
}

type Module struct {
	Name             string
	InstallProcess   Process
	UninstallProcess Process
}

var Nop = func(ctx context.Context, plugin *Plugin) {}

type Process func(ctx context.Context, plugin *Plugin, next func(ctx context.Context, plugin *Plugin))

func ComposeProcess(series ...Process) Process {
	return func(ctx context.Context, plugin *Plugin, receiver func(ctx context.Context, plugin *Plugin)) {
		if len(series) == 0 {
			receiver(ctx, plugin)
			return
		}
		series[0](ctx, plugin, func(newctx context.Context, plugin *Plugin) {
			ComposeProcess(series[1:]...)(newctx, plugin, receiver)
		})
	}
}

func CreateModule(name string, installfn Process, uninstallfn Process) *Module {
	return &Module{
		Name:             name,
		InstallProcess:   installfn,
		UninstallProcess: uninstallfn,
	}
}

type Options struct {
	Entry        string
	StartCommand string
	Modules      []*Module
}

func CreatePlugin(opt *Options) *Plugin {
	p := New()
	p.entry = opt.Entry
	p.startCommand = opt.StartCommand
	p.modules = opt.Modules
	return p
}
