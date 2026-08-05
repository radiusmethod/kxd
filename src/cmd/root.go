package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
)

// errConfigNotFound is returned when argv named a kubeconfig that isn't in
// ~/.kube. The generated shell function keys off the exit code, so this has to
// fail rather than exit 0.
var errConfigNotFound = errors.New("config does not exist")

var rootCmd = &cobra.Command{
	Use:   "kxd",
	Short: "kxd - switch between Kubeconfigs and contexts.",
	Long:  "Allows for switching kubeconfig files and contexts, as well as getting the current set ones.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runConfigSwitcher()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// RootCmd returns the root cobra command. Used by docs generation.
func RootCmd() *cobra.Command {
	return rootCmd
}

func Execute() {
	if shouldRunDirectConfigSwitch() {
		config := os.Args[1]
		handleSwitchError(directConfigSwitch(config))
		return
	}
	runRootCmd()
}

func directConfigSwitch(desiredConfig string) error {
	configs := utils.GetConfigs()
	if utils.Contains(configs, desiredConfig) {
		fmt.Printf(utils.PromptColor, "Config ")
		fmt.Printf(utils.CyanColor, desiredConfig)
		fmt.Printf(utils.PromptColor, " set.\n")
		if desiredConfig == "default" {
			desiredConfig = "config"
		}
		utils.WriteFile(desiredConfig, utils.GetHomeDir())
		return nil
	}

	// stderr, not stdout: the shell integration evals command substitutions of
	// this binary, so a warning on stdout would be executed instead of shown.
	fmt.Fprintf(os.Stderr, utils.NoticeColor, "WARNING: Config ")
	fmt.Fprintf(os.Stderr, utils.CyanColor, desiredConfig)
	fmt.Fprintf(os.Stderr, utils.NoticeColor, " does not exist or is invalid.\n")

	return fmt.Errorf("%w: %s", errConfigNotFound, desiredConfig)
}

// handleSwitchError exits non-zero when a config switch fails. A missing config
// has already been reported on stderr, so exit quietly instead of printing it a
// second time through log.Fatal.
func handleSwitchError(err error) {
	if err == nil {
		return
	}
	if errors.Is(err, errConfigNotFound) {
		os.Exit(1)
	}
	log.Fatal(err)
}

func runRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func shouldRunDirectConfigSwitch() bool {
	// Any argv[1] not in this list is treated as a kubeconfig name, so every
	// subcommand and alias has to be listed here.
	invalidConfigs := []string{"f", "file", "ctx", "context", "ns", "namespace", "init", "shellenv", "completion", "help", "--help", "v", "version"}
	return len(os.Args) > 1 && !utils.Contains(invalidConfigs, os.Args[1])
}
