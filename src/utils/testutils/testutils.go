package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempDir creates a temp dir for testing and returns its path.
func CreateTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "kxd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// CreateTempHome makes a temp dir, creates ~/.kube/ inside it, points $HOME at
// it, and returns the home path. This is the standard setup for kxd tests
// because GetConfigFileLocation log.Fatals if ~/.kube/ is missing.
func CreateTempHome(t *testing.T) string {
	t.Helper()
	home := CreateTempDir(t)
	t.Cleanup(func() { CleanupTempDir(t, home) })
	kubeDir := filepath.Join(home, ".kube")
	if err := os.MkdirAll(kubeDir, 0755); err != nil {
		t.Fatalf("Failed to create .kube dir: %v", err)
	}
	t.Setenv("HOME", home)
	return home
}

// WriteKubeconfigFixture copies a named fixture from testdata/kubeconfigs/ to
// dest and returns dest. The caller picks the destination (typically
// ~/.kube/config or ~/.kube/<name>.conf). repoRoot is the relative path from
// the test file to the repo root (e.g. "../.." from src/utils).
func WriteKubeconfigFixture(t *testing.T, repoRoot, fixture, dest string) string {
	t.Helper()
	src := filepath.Join(repoRoot, "testdata", "kubeconfigs", fixture)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read fixture %s: %v", src, err)
	}
	if err := os.WriteFile(dest, data, 0644); err != nil {
		t.Fatalf("Failed to write kubeconfig to %s: %v", dest, err)
	}
	return dest
}

// CleanupTempDir removes a temp dir and its contents.
func CleanupTempDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

// UnsetTestEnv removes an env var and checks the error. t.Setenv can't unset,
// so this is used when a test needs the variable explicitly absent.
func UnsetTestEnv(t *testing.T, key string) {
	t.Helper()
	if err := os.Unsetenv(key); err != nil {
		t.Errorf("Failed to unset environment variable %s: %v", key, err)
	}
}
