package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var contextCmd = &cobra.Command{
	Use:     "context",
	Short:   "Kubeconfig context command",
	Aliases: []string{"ctx"},
	Long:    "This is the default context command.",
}

var currentContextCmd = &cobra.Command{
	Use:     "current",
	Short:   "Shows currently set kubeconfig context",
	Aliases: []string{"c"},
	Long:    "This shows the current set kubeconfig context.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runListContext()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var switchContextCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch kubeconfig contexts",
	Aliases: []string{"s"},
	Long:    "This allows for switching of your kubeconfig context.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runContextSwitcher()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	contextCmd.AddCommand(switchContextCmd)
	contextCmd.AddCommand(currentContextCmd)
	rootCmd.AddCommand(contextCmd)
}

func runContextSwitcher() error {
	homeDir := utils.GetHomeDir()
	kubeconfigPath := utils.GetEnv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	config, err := utils.InitializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	contexts := utils.ListContexts(kubeconfigPath)

	fmt.Printf(utils.NoticeColor, "Kubeconfig Context Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a context"),
		Items:        contexts,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | magenta }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | magenta }}",
		},
		Searcher:          utils.NewPromptUISearcher(contexts),
		StartInSearchMode: true,
		Stdout:            &utils.BellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.CheckError(err)
	}

	fmt.Printf(utils.PromptColor, "Choose a context")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.MagentaColor, result)
	fmt.Println("")

	err = utils.SwitchContext(config, result, kubeconfigPath)
	if err != nil {
		log.Fatalf("Error switching context: %v\n", err)
	}
	return nil
}

func runListContext() error {
	homeDir := utils.GetHomeDir()
	kubeconfigPath := utils.GetEnv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	config, err := utils.InitializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}
	currentContext := config.CurrentContext
	if currentContext == "" {
		fmt.Println("No current context found in kubeconfig.")
		os.Exit(1)
	}
	fmt.Printf(currentContext + "\n")
	return nil
}
