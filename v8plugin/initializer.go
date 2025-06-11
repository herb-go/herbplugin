package v8plugin

import (
	"github.com/herb-go/herbplugin"
	v8 "github.com/herb-go/v8go"
)

type Initializer struct {
	Entry          string
	StartCommand   string
	DisableBuiltin bool
	Namespace      string
	Modules        []*herbplugin.Module
}

func (i *Initializer) MustApplyInitializer(p *Plugin) {
	p.Runtime = v8.NewContext()
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
