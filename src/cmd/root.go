package cmd

import (
	"fmt"
	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

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

func Execute() {
	if shouldRunDirectConfigSwitch() {
		config := os.Args[1]
		if err := directConfigSwitch(config); err != nil {
			log.Fatal(err)
		}
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

	fmt.Printf(utils.NoticeColor, "WARNING: Config ")
	fmt.Printf(utils.CyanColor, desiredConfig)
	fmt.Printf(utils.NoticeColor, " does not exist or is invalid.\n")

	return nil
}

func runRootCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func shouldRunDirectConfigSwitch() bool {
	invalidConfigs := []string{"f", "file", "ctx", "context", "completion", "help", "--help", "v", "version"}
	return len(os.Args) > 1 && !utils.Contains(invalidConfigs, os.Args[1])
}
