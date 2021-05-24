package luaplugin

import "testing"

func TestPlugin(t *testing.T) {
	p := newTestPlugin()
	if p.PluginType() != PluginType {
		t.Fatal(p)
	}
}
