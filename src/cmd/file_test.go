package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFileCommandMetadata(t *testing.T) {
	assert.Equal(t, "file", fileCmd.Use)
	assert.Equal(t, []string{"f"}, fileCmd.Aliases)
}

func TestRunConfigLister(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "dev.conf"), []byte("x"), 0644))

	// Should print without error. Output isn't asserted — runConfigLister
	// writes to stdout, and we already cover discovery in TestGetConfigs*.
	assert.NoError(t, runConfigLister())
}

func TestRunGetCurrentConfigFromKxd(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	// Selected kubeconfig exists at ~/.kube/dev.conf
	devPath := filepath.Join(kube, "dev.conf")
	assert.NoError(t, os.WriteFile(devPath, []byte("x"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(home, ".kxd"), []byte("dev.conf"), 0644))

	assert.NoError(t, runGetCurrentConfig())
}

func TestRunGetCurrentConfigKxdMissingFile(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	// .kxd points at a non-existent file → falls back to default ~/.kube/config,
	// which we make sure exists.
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "config"), []byte("x"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(home, ".kxd"), []byte("ghost.conf"), 0644))

	assert.NoError(t, runGetCurrentConfig())
}

func TestRunGetCurrentConfigNoKxdUsesDefault(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "config"), []byte("x"), 0644))
	// No .kxd file, no KUBECONFIG → uses ~/.kube/config.
	testutils.UnsetTestEnv(t, "KUBECONFIG")

	assert.NoError(t, runGetCurrentConfig())
}

func TestRunGetCurrentConfigKubeconfigEnv(t *testing.T) {
	home := testutils.CreateTempHome(t)
	explicit := filepath.Join(home, "elsewhere.conf")
	assert.NoError(t, os.WriteFile(explicit, []byte("x"), 0644))
	t.Setenv("KUBECONFIG", explicit)

	assert.NoError(t, runGetCurrentConfig())
}
