package lua51plugin

import (
	"testing"

	"github.com/herb-go/herbplugin"
	lua "github.com/yuin/gopher-lua"
)

var output = ""
var testInitializer = NewInitializer()

func init() {
	testInitializer.Entry = "main.lua"
	testInitializer.StartCommand = "start()"
	AppendCommonModules(testInitializer)
}

type testPlugin struct {
	*Plugin
}

func (p *testPlugin) MustLoadParam(name string) string {
	p.Lock()
	defer p.Unlock()
	if err := p.LState.CallByParam(lua.P{
		Fn:      p.LState.GetGlobal("getparam"),
		NRet:    1,
		Protect: true,
	}, lua.LString(name)); err != nil {
		panic(err)
	}
	ret := p.LState.Get(-1).(lua.LString)
	p.LState.Pop(1)
	return ret.String()
}

func newTestPlugin() *testPlugin {
	return &testPlugin{
		Plugin: MustCreatePlugin(testInitializer),
	}
}
func TestCommonModule(t *testing.T) {
	o := herbplugin.NewOptions()
	o.Location.Path = "testscripts"
	o.Params["testkey"] = "testvalue"
	output = ""
	p := newTestPlugin()
	p.SetPluginPrinter(func(info string) {
		output = output + info
	})
	herbplugin.Lanuch(p, o)
	if output != "printed" {
		t.Fatal(output)
	}
	defer p.MustClosePlugin()
	param := p.MustLoadParam("test")
	if param != "" {
		t.Fatal()
	}
	param = p.MustLoadParam("testkey")
	if param != "testvalue" {
		t.Fatal(param)
	}

}
