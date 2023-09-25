package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "kxd",
	Short: "kxd - switch between Kubeconfigs and contexts.",
	Long:  "Allows for switching kubeconfig files and contexts, as well as getting the current set ones.",
	// Should we default to kxd file switch?
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
