package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShouldRunDirectConfigSwitch(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{"direct config name", []string{"kxd", "dev.conf"}, true},
		{"file subcommand long", []string{"kxd", "file"}, false},
		{"file subcommand short", []string{"kxd", "f"}, false},
		{"context subcommand long", []string{"kxd", "context"}, false},
		{"context subcommand short", []string{"kxd", "ctx"}, false},
		{"namespace subcommand long", []string{"kxd", "namespace"}, false},
		{"namespace subcommand short", []string{"kxd", "ns"}, false},
		{"completion", []string{"kxd", "completion"}, false},
		{"help long", []string{"kxd", "--help"}, false},
		{"help short", []string{"kxd", "help"}, false},
		{"version long", []string{"kxd", "version"}, false},
		{"version short", []string{"kxd", "v"}, false},
		{"no args", []string{"kxd"}, false},
	}

	original := os.Args
	defer func() { os.Args = original }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			assert.Equal(t, tt.expected, shouldRunDirectConfigSwitch())
		})
	}
}

func TestDirectConfigSwitch(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	// Configs available: dev.conf, prod.conf, plus the default "config" file
	// which makes the "default" sentinel valid.
	for _, name := range []string{"dev.conf", "prod.conf", "config"} {
		assert.NoError(t, os.WriteFile(filepath.Join(kube, name), []byte("x"), 0644))
	}

	tests := []struct {
		name         string
		desired      string
		expectFile   bool
		expectedKxd  string
	}{
		{"valid config", "dev.conf", true, "dev.conf"},
		// "default" is rewritten to "config" before being persisted.
		{"default sentinel rewrites to config", "default", true, "config"},
		{"invalid config leaves .kxd untouched", "nope.conf", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kxdFile := filepath.Join(home, ".kxd")
			_ = os.Remove(kxdFile)

			err := directConfigSwitch(tt.desired)
			assert.NoError(t, err)

			if tt.expectFile {
				content, err := os.ReadFile(kxdFile)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedKxd, string(content))
			} else {
				_, err := os.Stat(kxdFile)
				assert.True(t, os.IsNotExist(err), ".kxd should not exist for invalid config")
			}
		})
	}
}

func TestRootCmdMetadata(t *testing.T) {
	cmd := RootCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "kxd", cmd.Use)
	assert.Contains(t, cmd.Short, "kxd")
}
