package cmd

import (
	"path/filepath"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestContextCommandMetadata(t *testing.T) {
	assert.Equal(t, "context", contextCmd.Use)
	assert.Equal(t, []string{"ctx"}, contextCmd.Aliases)
}

func TestRunListContext(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	testutils.WriteKubeconfigFixture(t, "../..", "basic.conf", cfgPath)
	t.Setenv("KUBECONFIG", cfgPath)

	assert.NoError(t, runListContext())
}

func TestRunContextLister(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	testutils.WriteKubeconfigFixture(t, "../..", "basic.conf", cfgPath)
	t.Setenv("KUBECONFIG", cfgPath)

	assert.NoError(t, runContextLister())
}
