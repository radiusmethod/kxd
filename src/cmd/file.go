package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/radiusmethod/kxd/src/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var fileCmd = &cobra.Command{
	Use:     "file",
	Short:   "Kubeconfig file command",
	Aliases: []string{"f"},
	Long:    "This is the default file command.",
}

var currentFileCmd = &cobra.Command{
	Use:     "current",
	Short:   "Shows currently set kubeconfig",
	Aliases: []string{"c"},
	Long:    "This shows the current set kubeconfig file.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runGetCurrentConfig()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var switchFileCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Switch kubeconfig",
	Aliases: []string{"s"},
	Long:    "This allows for switching of your kubeconfig.",
	Run: func(cmd *cobra.Command, args []string) {
		err := runConfigSwitcher()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	fileCmd.AddCommand(switchFileCmd)
	fileCmd.AddCommand(currentFileCmd)
	rootCmd.AddCommand(fileCmd)
}

func runConfigSwitcher() error {
	homeDir := utils.GetHomeDir()
	configFileLocation := fmt.Sprintf("%s/.kube", homeDir)
	configs := getConfigs(configFileLocation, homeDir)

	err := utils.TouchFile(fmt.Sprintf("%s/.kxd", homeDir))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(utils.NoticeColor, "Kubeconfig Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(utils.PromptColor, "Choose a config"),
		Items:        configs,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Searcher:          utils.NewPromptUISearcher(configs),
		StartInSearchMode: true,
		Stdout:            &utils.BellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		utils.CheckError(err)
	}

	fmt.Printf(utils.PromptColor, "Choose a config")
	fmt.Printf(utils.NoticeColor, "? ")
	fmt.Printf(utils.CyanColor, result)
	fmt.Println("")

	if result == "default" {
		result = "config"
	}
	utils.WriteFile(result, homeDir)

	return nil
}

func runGetCurrentConfig() error {
	homeDir := utils.GetHomeDir()
	kubeconfigPath := utils.GetEnv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	if _, err := os.Stat(kubeconfigPath); err == nil {
		fmt.Println(kubeconfigPath)
	} else {
		fmt.Println("No current kubeconfig found.")
		os.Exit(1)
	}
	return nil
}

func getConfigs(configFileLocation string, homeDir string) []string {
	var files []string
	fileExts := strings.Split(utils.GetEnv("KXD_MATCHER", ".conf"), ",")
	err := filepath.Walk(configFileLocation, func(path string, f os.FileInfo, _ error) error {
		for _, value := range fileExts {
			if !f.IsDir() && strings.Contains(f.Name(), value) {
				files = append(files, f.Name())
				break
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/.kube/config", homeDir)); err == nil {
		files = append(files, "default")
	}
	files = append(files, "unset")
	sort.Strings(files)
	return files
}
