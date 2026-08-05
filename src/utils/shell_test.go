package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseShell(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Shell
		expectError bool
	}{
		{name: "bash", input: "bash", expected: Bash},
		{name: "zsh", input: "zsh", expected: Zsh},
		{name: "fish", input: "fish", expected: Fish},
		{name: "powershell", input: "powershell", expected: PowerShell},
		{name: "pwsh alias", input: "pwsh", expected: PowerShell},
		{name: "case insensitive", input: "PowerShell", expected: PowerShell},
		{name: "surrounding space", input: "  zsh ", expected: Zsh},
		{name: "unknown shell", input: "csh", expectError: true},
		{name: "empty", input: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := ParseShell(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, shell)
		})
	}
}

func TestReadState(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		write    bool
		expected string
	}{
		{name: "Config name", contents: "dev.conf\n", write: true, expected: "dev.conf"},
		{name: "Default config", contents: "config\n", write: true, expected: "config"},
		{name: "Empty means unset", contents: "", write: true, expected: ""},
		{name: "Trailing whitespace trimmed", contents: "  dev.conf  \n", write: true, expected: "dev.conf"},
		{name: "Missing file reads as empty", write: false, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.write {
				if err := os.WriteFile(filepath.Join(dir, ".kxd"), []byte(tt.contents), 0644); err != nil {
					t.Fatalf("Failed to write .kxd: %v", err)
				}
			}
			got, err := ReadState(dir)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestShellEnv(t *testing.T) {
	home := "/home/tester"

	tests := []struct {
		name     string
		config   string
		shell    Shell
		expected string
	}{
		{
			name:     "Named config, posix",
			config:   "dev.conf",
			shell:    Bash,
			expected: "export KUBECONFIG='/home/tester/.kube/dev.conf'\n",
		},
		{
			name:     "Default config maps to ~/.kube/config",
			config:   "config",
			shell:    Zsh,
			expected: "export KUBECONFIG='/home/tester/.kube/config'\n",
		},
		{
			name:     "Empty unsets, posix",
			config:   "",
			shell:    Bash,
			expected: "unset KUBECONFIG\n",
		},
		{
			name:     "Named config, fish",
			config:   "dev.conf",
			shell:    Fish,
			expected: "set -gx KUBECONFIG '/home/tester/.kube/dev.conf'\n",
		},
		{
			name:     "Empty unsets, fish",
			config:   "",
			shell:    Fish,
			expected: "set -e KUBECONFIG\n",
		},
		{
			name:     "Named config, powershell",
			config:   "dev.conf",
			shell:    PowerShell,
			expected: "$env:KUBECONFIG = '/home/tester/.kube/dev.conf'\n",
		},
		{
			name:     "Empty unsets, powershell",
			config:   "",
			shell:    PowerShell,
			expected: "$env:KUBECONFIG = $null\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ShellEnv(tt.config, home, tt.shell))
		})
	}
}

// Kubeconfig filenames come off the filesystem, so they can contain characters
// the shell would otherwise interpret. Unquoted output would run as code.
func TestShellEnvQuoting(t *testing.T) {
	home := "/home/tester"

	tests := []struct {
		name     string
		shell    Shell
		config   string
		expected string
	}{
		{
			name:     "Space, posix",
			shell:    Bash,
			config:   "my cluster.conf",
			expected: "export KUBECONFIG='/home/tester/.kube/my cluster.conf'\n",
		},
		{
			name:     "Single quote, posix",
			shell:    Bash,
			config:   "we'ird.conf",
			expected: `export KUBECONFIG='/home/tester/.kube/we'\''ird.conf'` + "\n",
		},
		{
			name:     "Command substitution stays literal, posix",
			shell:    Bash,
			config:   "$(whoami).conf",
			expected: "export KUBECONFIG='/home/tester/.kube/$(whoami).conf'\n",
		},
		{
			name:     "Single quote, fish",
			shell:    Fish,
			config:   "we'ird.conf",
			expected: `set -gx KUBECONFIG '/home/tester/.kube/we\'ird.conf'` + "\n",
		},
		{
			name:     "Single quote, powershell",
			shell:    PowerShell,
			config:   "we'ird.conf",
			expected: "$env:KUBECONFIG = '/home/tester/.kube/we''ird.conf'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ShellEnv(tt.config, home, tt.shell))
		})
	}
}

func TestInitScript(t *testing.T) {
	tests := []struct {
		name     string
		shell    Shell
		contains []string
	}{
		{
			name:     "bash",
			shell:    Bash,
			contains: []string{"kxd() {", `eval "$(command kxd shellenv bash)"`, "complete -o nospace -F _kxd_completion kxd"},
		},
		{
			name:     "zsh",
			shell:    Zsh,
			contains: []string{"kxd() {", `eval "$(command kxd shellenv zsh)"`, "bashcompinit"},
		},
		{
			name:     "fish",
			shell:    Fish,
			contains: []string{"function kxd", "command kxd shellenv fish | source", "complete -c kxd"},
		},
		{
			name:     "powershell",
			shell:    PowerShell,
			contains: []string{"function kxd", "& $global:KxdBin shellenv powershell | Out-String | Invoke-Expression", "Register-ArgumentCompleter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := InitScript(tt.shell)
			assert.NotEmpty(t, script)
			for _, want := range tt.contains {
				assert.Contains(t, script, want)
			}
		})
	}
}

// `kxd namespace list` calls the live cluster, so wiring it into completion
// would hang the shell on a TAB press.
func TestInitScriptDoesNotCompleteNamespaces(t *testing.T) {
	for _, shell := range []Shell{Bash, Zsh, Fish, PowerShell} {
		assert.NotContains(t, InitScript(shell), "namespace list", "%s must not shell out to namespace list", shell)
	}
}

// bashcompinit only provides `complete` for zsh; bash must not carry it.
func TestInitScriptBashHasNoCompinit(t *testing.T) {
	assert.NotContains(t, InitScript(Bash), "bashcompinit")
}

func TestInitScriptUnknownShell(t *testing.T) {
	assert.Empty(t, InitScript(Shell("csh")))
}

// Every accepted shell has to produce a script, or `kxd init <shell>` would
// validate the argument and then print nothing.
func TestInitScriptCoversAcceptedShells(t *testing.T) {
	for _, name := range AcceptedShells {
		shell, err := ParseShell(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, InitScript(shell), "no init script for %s", name)
		assert.True(t, strings.HasSuffix(InitScript(shell), "\n"), "%s script must end in a newline", name)
	}
}
