package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"log"
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
