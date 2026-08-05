package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Print the shell integration for the given shell.",
	Long: "Print the shell integration for the given shell: the kxd function, completion, " +
		"and the hook that applies the kubeconfig selected in ~/.kxd to new shells.\n\n" +
		"kxd has to run inside your shell to export KUBECONFIG, because a child process " +
		"cannot change its parent's environment. Eval this from your rc file:\n\n" +
		"  bash/zsh:   eval \"$(kxd init zsh)\"\n" +
		"  fish:       kxd init fish | source\n" +
		"  PowerShell: kxd init powershell | Out-String | Invoke-Expression",
	ValidArgs: utils.AcceptedShells,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := utils.ParseShell(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(utils.InitScript(shell))
	},
}

var shellenvCmd = &cobra.Command{
	Use:   "shellenv [shell]",
	Short: "Print the shell code that applies ~/.kxd to the current shell.",
	Long: "Print the export/unset statement for the selected kubeconfig. " +
		"Used by the function that `kxd init` generates; you should not need to call it directly. " +
		"Defaults to POSIX (bash/zsh) syntax.",
	ValidArgs: utils.AcceptedShells,
	Args:      cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		shell := utils.Bash
		if len(args) == 1 {
			s, err := utils.ParseShell(args[0])
			if err != nil {
				log.Fatal(err)
			}
			shell = s
		}
		homeDir := utils.GetHomeDir()
		config, err := utils.ReadState(homeDir)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(utils.ShellEnv(config, homeDir, shell))
	},
}

func init() {
	initCmd.Example = fmt.Sprintf("  eval \"$(kxd init zsh)\"\n\n  # supported: %s", strings.Join(utils.SupportedShells, ", "))
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(shellenvCmd)
}
