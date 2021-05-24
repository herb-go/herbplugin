package luaplugin

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

const PluginType = "lua"
const DefaultNamespace = "plugin"

func New() *Plugin {
	return &Plugin{
		BasicPlugin: herbplugin.New(),
	}
}

type Plugin struct {
	sync.Mutex
	entry  string
	LState *lua.LState
	*herbplugin.BasicPlugin
	DisableBuiltin bool
	startCommand   string
	modules        []*Module
	namespace      string
	Builtin        map[string]lua.LGFunction
}

func (p *Plugin) PluginType() string {
	return PluginType
}
func (p *Plugin) MustInitPlugin() {
	p.BasicPlugin.MustInitPlugin()
	p.Builtin = map[string]lua.LGFunction{}
	var processs = make([]Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].InitProcess != nil {
			processs = append(processs, p.modules[k].InitProcess)
		}
	}
	ComposeProcess(processs...)(context.TODO(), p, Nop)
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
	var processs = make([]Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].BootProcess != nil {
			processs = append(processs, p.modules[k].BootProcess)
		}
	}
	ComposeProcess(processs...)(context.TODO(), p, Nop)
	if p.startCommand != "" {
		err := p.LState.DoString(p.startCommand)
		if err != nil {
			panic(err)
		}
	}
}

func (p *Plugin) MustClosePlugin() {
	defer p.LState.Close()
	var processs = make([]Process, 0, len(p.modules))
	for k := range p.modules {
		if p.modules[k].CloseProcess != nil {
			processs = append(processs, p.modules[k].CloseProcess)
		}
	}
	ComposeProcess(processs...)(context.TODO(), p, Nop)
	p.BasicPlugin.MustClosePlugin()
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

type Initializer struct {
	Entry          string
	StartCommand   string
	DisableBuiltin bool
	Namespace      string
	Modules        []*Module
	Options        []lua.Options
}

func (i *Initializer) MustApplyInitializer(p *Plugin) {
	p.LState = lua.NewState(i.Options...)
	p.entry = i.Entry
	p.startCommand = i.StartCommand
	p.modules = i.Modules
	p.namespace = i.Namespace
	if p.namespace == "" {
		p.namespace = DefaultNamespace
	}
	p.DisableBuiltin = i.DisableBuiltin
}

func NewInitializer() *Initializer {
	return &Initializer{}
}
func MustCreatePlugin(i *Initializer) *Plugin {
	p := New()
	i.MustApplyInitializer(p)
	return p
}
