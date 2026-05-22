package utils

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func writeBasicKubeconfig(t *testing.T, path string) {
	t.Helper()
	testutils.WriteKubeconfigFixture(t, "../..", "basic.conf", path)
}

func TestInitializeKubeconfig(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	writeBasicKubeconfig(t, cfgPath)

	cfg, err := InitializeKubeconfig(cfgPath)
	assert.NoError(t, err)
	assert.Equal(t, "context-a", cfg.CurrentContext)
	assert.Len(t, cfg.Contexts, 2)
}

func TestInitializeKubeconfigMissing(t *testing.T) {
	_, err := InitializeKubeconfig("/nonexistent/kubeconfig")
	assert.Error(t, err)
}

func TestListContexts(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	writeBasicKubeconfig(t, cfgPath)

	contexts := ListContexts(cfgPath)
	assert.Equal(t, []string{"context-a", "context-b"}, contexts)
	assert.True(t, sort.StringsAreSorted(contexts))
}

func TestSwitchContext(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	writeBasicKubeconfig(t, cfgPath)

	cfg, err := InitializeKubeconfig(cfgPath)
	assert.NoError(t, err)

	assert.NoError(t, SwitchContext(cfg, "context-b", cfgPath))

	reread, err := InitializeKubeconfig(cfgPath)
	assert.NoError(t, err)
	assert.Equal(t, "context-b", reread.CurrentContext)
}

func TestSwitchNamespace(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	writeBasicKubeconfig(t, cfgPath)

	assert.NoError(t, SwitchNamespace("kube-system", cfgPath))

	cfg, err := InitializeKubeconfig(cfgPath)
	assert.NoError(t, err)
	assert.Equal(t, "kube-system", cfg.Contexts["context-a"].Namespace)
}

func TestGetConfigsDefaultMatcher(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	for _, name := range []string{"dev.conf", "prod.conf", "notes.txt"} {
		assert.NoError(t, os.WriteFile(filepath.Join(kube, name), []byte("x"), 0644))
	}

	configs := GetConfigs()
	assert.Contains(t, configs, "dev.conf")
	assert.Contains(t, configs, "prod.conf")
	assert.NotContains(t, configs, "notes.txt")
	// "unset" sentinel is always present; "default" only appears when
	// ~/.kube/config exists, which it doesn't here.
	assert.Contains(t, configs, "unset")
	assert.NotContains(t, configs, "default")
}

func TestGetConfigsWithDefault(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	assert.NoError(t, os.WriteFile(filepath.Join(kube, "config"), []byte("x"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "dev.conf"), []byte("x"), 0644))

	configs := GetConfigs()
	assert.Contains(t, configs, "default", "default sentinel should appear when ~/.kube/config exists")
	assert.Contains(t, configs, "dev.conf")
	assert.Contains(t, configs, "unset")
}

func TestGetConfigsCustomMatcher(t *testing.T) {
	home := testutils.CreateTempHome(t)
	kube := filepath.Join(home, ".kube")

	assert.NoError(t, os.WriteFile(filepath.Join(kube, "team-config"), []byte("x"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "prod.conf"), []byte("x"), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(kube, "readme.md"), []byte("x"), 0644))

	t.Setenv("KXD_MATCHER", "-config,.conf")

	configs := GetConfigs()
	assert.Contains(t, configs, "team-config", "Custom matcher should pick up '-config' suffix")
	assert.Contains(t, configs, "prod.conf")
	assert.NotContains(t, configs, "readme.md")
}
