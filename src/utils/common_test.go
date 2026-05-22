package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/radiusmethod/kxd/src/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTouchFile(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	err := TouchFile(filePath)
	assert.NoError(t, err)

	_, err = os.Stat(filePath)
	assert.NoError(t, err, "File should exist")
}

func TestWriteFile(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	WriteFile("dev.conf", tempDir)
	filePath := filepath.Join(tempDir, ".kxd")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "dev.conf", string(content), "Should write config name")

	// "unset" is the sentinel that empties the file → wrapper unsets KUBECONFIG.
	WriteFile("unset", tempDir)
	content, err = os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "", string(content), "'unset' should write empty content")
}

func TestGetEnv(t *testing.T) {
	t.Setenv("TEST_KXD_VAR", "test-value")
	assert.Equal(t, "test-value", GetEnv("TEST_KXD_VAR", "fallback"))
	assert.Equal(t, "fallback", GetEnv("NONEXISTENT_KXD_VAR", "fallback"))
}

func TestGetHomeDir(t *testing.T) {
	t.Setenv("HOME", "/test/home")
	assert.Equal(t, "/test/home", GetHomeDir())
}

func TestGetConfigFileLocation(t *testing.T) {
	home := testutils.CreateTempHome(t)

	expected := filepath.Join(home, ".kube")
	assert.Equal(t, expected, GetConfigFileLocation())
}

func TestGetCurrentConfigFile(t *testing.T) {
	t.Setenv("KUBECONFIG", "/explicit/config")
	assert.Equal(t, "/explicit/config", GetCurrentConfigFile())

	testutils.UnsetTestEnv(t, "KUBECONFIG")
	t.Setenv("HOME", "/somewhere")
	assert.Equal(t, "/somewhere/.kube/config", GetCurrentConfigFile())
}

func TestIsDirectoryExists(t *testing.T) {
	tempDir := testutils.CreateTempDir(t)
	defer testutils.CleanupTempDir(t, tempDir)

	assert.True(t, IsDirectoryExists(tempDir))
	assert.False(t, IsDirectoryExists("/definitely/not/a/real/dir/12345"))

	// Files are not directories.
	filePath := filepath.Join(tempDir, "afile")
	assert.NoError(t, os.WriteFile(filePath, []byte("x"), 0644))
	assert.False(t, IsDirectoryExists(filePath))
}

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	assert.True(t, Contains(slice, "b"))
	assert.False(t, Contains(slice, "d"))

	var empty []string
	assert.False(t, Contains(empty, "a"))
}

// CheckError fatally exits, so each branch runs in a subprocess and we verify
// the exit status.
func TestCheckErrorDel(t *testing.T) {
	if os.Getenv("TEST_KXD_CHECK_ERROR_DEL") == "1" {
		CheckError(fmt.Errorf("^D"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErrorDel")
	cmd.Env = append(os.Environ(), "TEST_KXD_CHECK_ERROR_DEL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestCheckErrorCtrlC(t *testing.T) {
	if os.Getenv("TEST_KXD_CHECK_ERROR_CTRL_C") == "1" {
		CheckError(fmt.Errorf("^C"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErrorCtrlC")
	cmd.Env = append(os.Environ(), "TEST_KXD_CHECK_ERROR_CTRL_C=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}

func TestCheckErrorOther(t *testing.T) {
	if os.Getenv("TEST_KXD_CHECK_ERROR_OTHER") == "1" {
		CheckError(fmt.Errorf("other error"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErrorOther")
	cmd.Env = append(os.Environ(), "TEST_KXD_CHECK_ERROR_OTHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}
