package cmd

import (
	"fmt"
	"github.com/radiusmethod/kxd/src/utils"
	"github.com/radiusmethod/promptui"
	"github.com/spf13/cobra"
	"log"
)

var namespaceCmd = &cobra.Command{
	Use:     "namespace",
	Short:   "Kubeconfig namespace command",
	Aliases: []string{"ns"},
	Long:    "This is the default namespace command.",
}

var currentNamespaceCmd = &cobra.Command{
	Use:     "current",
	Short:   "Shows currently set k8s namespace",
	Aliases: []string{"c"},
	Long:    "This shows the current set k8s namespace.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runListNamespace()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var switchNamespaceCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch k8s namespaces",
	Aliases: []string{"s"},
	Long:    "This allows for switching of your k8s namespace.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runNamespaceSwitcher()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var listNamespacesCmd = &cobra.Command{
	Use:     "list",
	Short:   "List kubeconfig namespaces",
	Aliases: []string{"l"},
	Long:    "This displays a simple list of your k8s namespaces.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runNamespaceLister()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	namespaceCmd.AddCommand(switchNamespaceCmd, currentNamespaceCmd, listNamespacesCmd)
	rootCmd.AddCommand(namespaceCmd)
}

func runNamespaceSwitcher() error {
	kubeconfigPath := utils.GetCurrentConfigFile()
	namespaces := utils.ListNamespaces(kubeconfigPath)

	fmt.Printf(utils.NoticeColor, "Kubeconfig Namespace Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a namespace"),
		Items:        namespaces,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | magenta }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | magenta }}",
		},
		Searcher:          utils.NewPromptUISearcher(namespaces),
		StartInSearchMode: true,
		Stdout:            &utils.BellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.CheckError(err)
	}

	fmt.Printf(utils.PromptColor, "Choose a namespace")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.MagentaColor, result)
	fmt.Println("")

	err = utils.SwitchNamespace(result, kubeconfigPath)
	if err != nil {
		log.Fatalf("Error switching namespace: %v\n", err)
	}
	return nil
}

func runListNamespace() error {
	kubeconfigPath := utils.GetCurrentConfigFile()
	config, err := utils.InitializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}
	currentContext := config.CurrentContext
	context, exists := config.Contexts[currentContext]
	if !exists {
		log.Fatal("Current context not found in kubeconfig")
	}

	namespace := context.Namespace
	if namespace == "" {
		namespace = "default"
	}
	fmt.Printf(namespace + "\n")
	return nil
}

func runNamespaceLister() error {
	namespaces := utils.ListNamespaces(utils.GetCurrentConfigFile())
	for _, c := range namespaces {
		fmt.Println(c)
	}
	return nil
}
