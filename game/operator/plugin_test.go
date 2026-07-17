package operator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPluginLoadsAbsoluteUnixStylePath(t *testing.T) {
	pluginPath := filepath.Join(t.TempDir(), "absolute_plugin.py")
	if err := os.WriteFile(pluginPath, []byte("def Init():\n    pass\n"), 0600); err != nil {
		t.Fatal(err)
	}

	plugin := NewPlugin(pluginPath)
	defer plugin.ctx.Close()
	if err := plugin.LoadError(); err != nil {
		t.Fatalf("absolute plugin path was not loaded: %v", err)
	}
	if plugin.module == nil {
		t.Fatal("absolute plugin path produced a nil module")
	}
}

func TestNewPluginReportsMissingFile(t *testing.T) {
	plugin := NewPlugin(filepath.Join(t.TempDir(), "missing.py"))
	defer plugin.ctx.Close()
	if err := plugin.LoadError(); err == nil {
		t.Fatal("missing plugin did not report a load error")
	}
}
