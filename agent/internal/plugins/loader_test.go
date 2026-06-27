package plugins

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadExternalPlugins(t *testing.T) {
	root := t.TempDir()
	pluginDir := filepath.Join(root, "echo")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		t.Fatal(err)
	}

	executable, args, scriptPath := testPluginCommand(pluginDir)
	if scriptPath != "" {
		if err := os.WriteFile(scriptPath, []byte(testPluginScript()), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	manifest := `{
  "id": "external-echo",
  "name": "External Echo",
  "category": "utility",
  "description": "Test plugin",
  "commands": ["external_echo"],
  "executable": "` + executable + `",
  "args": [` + testArgsJSON(args) + `]
}`
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadExternalPlugins(root, log.New(os.Stdout, "", 0))
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(loaded))
	}

	registry := NewRegistry(loaded...)
	err = registry.Execute(context.Background(), Command{
		Type: "external",
		Name: "external_echo",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testPluginCommand(pluginDir string) (string, []string, string) {
	if runtime.GOOS == "windows" {
		scriptPath := filepath.Join(pluginDir, "plugin.cmd")
		return "plugin.cmd", nil, scriptPath
	}

	scriptPath := filepath.Join(pluginDir, "plugin.sh")
	return "plugin.sh", nil, scriptPath
}

func testPluginScript() string {
	if runtime.GOOS == "windows" {
		return "@echo off\r\nexit /b 0\r\n"
	}

	return "#!/bin/sh\ncat >/dev/null\nexit 0\n"
}

func testArgsJSON(args []string) string {
	result := ""
	for index, arg := range args {
		if index > 0 {
			result += ", "
		}
		result += `"` + filepath.ToSlash(arg) + `"`
	}
	return result
}
