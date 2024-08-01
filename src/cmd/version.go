package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version string = "v0.1.1"

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "kxd version command",
	Aliases: []string{"v"},
	Long:    "Returns the current version of kxd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kxd version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
