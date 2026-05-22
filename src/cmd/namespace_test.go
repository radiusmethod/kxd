package cmd

import (
	"path/filepath"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceCommandMetadata(t *testing.T) {
	assert.Equal(t, "namespace", namespaceCmd.Use)
	assert.Equal(t, []string{"ns"}, namespaceCmd.Aliases)
}

func TestRunListNamespaceWithExplicitNamespace(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	testutils.WriteKubeconfigFixture(t, "../..", "with_namespace.conf", cfgPath)
	t.Setenv("KUBECONFIG", cfgPath)

	// Current context has namespace=my-ns; should print it.
	assert.NoError(t, runListNamespace())
}

func TestRunListNamespaceDefaultsToDefault(t *testing.T) {
	home := testutils.CreateTempHome(t)
	cfgPath := filepath.Join(home, ".kube", "config")
	testutils.WriteKubeconfigFixture(t, "../..", "basic.conf", cfgPath)
	t.Setenv("KUBECONFIG", cfgPath)

	// basic.conf's current context has no namespace; should fall back to "default".
	assert.NoError(t, runListNamespace())
}

// runNamespaceLister and runNamespaceSwitcher hit a live cluster API, so they
// are not unit-testable without a kube fake. Covered by manual testing only.
