package luaplugin

import (
	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

type Initializer struct {
	Entry          string
	StartCommand   string
	DisableBuiltin bool
	Namespace      string
	Modules        []*herbplugin.Module
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
