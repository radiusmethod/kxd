package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	NoticeColor  = "\033[0;38m%s\u001B[0m"
	PromptColor  = "\033[1;38m%s\u001B[0m"
	CyanColor    = "\033[0;36m%s\033[0m"
	MagentaColor = "\033[0;35m%s\033[0m"
)

var version string = "v0.0.4"

func newPromptUISearcher(items []string) list.Searcher {
	return func(searchInput string, itemIndex int) bool {
		return strings.Contains(strings.ToLower(items[itemIndex]), strings.ToLower(searchInput))
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v\n", err)
	}

	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch arg {
		case "-v", "--v", "version":
			fmt.Println("kxd version:", version)
		case "-c", "--c", "context":
			err := runContextSwitcher(homeDir)
			if err != nil {
				log.Fatal(err)
			}
		case "-l", "--l", "list":
			err := runListConfig(homeDir)
			if err != nil {
				log.Fatal(err)
			}
		case "-lc", "--lc", "list-context":
			err := runListContext(homeDir)
			if err != nil {
				log.Fatal(err)
			}
		case "-h", "--h", "help":
			err := displayHelp()
			if err != nil {
				log.Fatal(err)
			}
		case "-s", "--s", "switch":
			err := runConfigSwitcher(homeDir)
			if err != nil {
				log.Fatal(err)
			}
		default:
			err := runConfigSwitcher(homeDir)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		err := runConfigSwitcher(homeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func displayHelp() error {
	var helpMessage strings.Builder
	options := map[string]string{
		"      [-s]  ": "Switch configs.",
		"       -c   ": "Switch contexts.",
		"       -l   ": "List current config.",
		"       -lc  ": "List current context.",
		"       -h   ": "Help. Displays this message.",
		"       -v   ": "Displays version.",
	}
	helpMessage.WriteString("Usage: kxd [OPERATION]\n")
	for option, description := range options {
		helpMessage.WriteString(fmt.Sprintf("  %s: %s\n", option, description))
	}
	fmt.Println(helpMessage.String())
	return nil
}

func runListConfig(homeDir string) error {
	kubeconfigPath := getenv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	if _, err := os.Stat(kubeconfigPath); err == nil {
		fmt.Println(kubeconfigPath)
	} else {
		fmt.Println("No current kubeconfig found.")
		os.Exit(1)
	}
	return nil
}

func runListContext(homeDir string) error {
	kubeconfigPath := getenv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	config, err := initializeKubeconfig(kubeconfigPath)
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

func runConfigSwitcher(homeDir string) error {
	configFileLocation := fmt.Sprintf("%s/.kube", homeDir)
	configs := getConfigs(configFileLocation, homeDir)

	err := touchFile(fmt.Sprintf("%s/.kxd", homeDir))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(NoticeColor, "Kubeconfig Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(PromptColor, "Choose a config"),
		Items:        configs,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Searcher:          newPromptUISearcher(configs),
		StartInSearchMode: true,
		Stdout:            &bellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		checkError(err)
	}

	fmt.Printf(PromptColor, "Choose a config")
	fmt.Printf(NoticeColor, "? ")
	fmt.Printf(CyanColor, result)
	fmt.Println("")

	if result == "default" {
		result = "config"
	}
	writeFile(result, homeDir)

	return nil
}

func runContextSwitcher(homeDir string) error {
	kubeconfigPath := getenv("KUBECONFIG", filepath.Join(homeDir, ".kube/config"))
	config, err := initializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	contexts := listContexts(kubeconfigPath)

	fmt.Printf(NoticeColor, "Kubeconfig Context Switcher\n")
	prompt := promptui.Select{
		Label:        fmt.Sprintf(PromptColor, "Choose a context"),
		Items:        contexts,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | magenta }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | magenta }}",
		},
		Searcher:          newPromptUISearcher(contexts),
		StartInSearchMode: true,
		Stdout:            &bellSkipper{},
	}

	_, result, err := prompt.Run()
	if err != nil {
		checkError(err)
	}

	fmt.Printf(PromptColor, "Choose a context")
	fmt.Printf(NoticeColor, "? ")
	fmt.Printf(MagentaColor, result)
	fmt.Println("")

	err = switchContext(config, result, kubeconfigPath)
	if err != nil {
		log.Fatalf("Error switching context: %v\n", err)
	}
	return nil
}

func initializeKubeconfig(kubeconfigPath string) (*api.Config, error) {
	return clientcmd.LoadFromFile(kubeconfigPath)
}

func listContexts(kubeconfigPath string) []string {
	config, err := initializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	var contexts []string
	for contextName := range config.Contexts {
		contexts = append(contexts, contextName)
	}
	return contexts
}

func switchContext(config *api.Config, contextName string, kubeconfigPath string) error {
	config.CurrentContext = contextName
	return clientcmd.WriteToFile(*config, kubeconfigPath)
}

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func writeFile(config, loc string) {
	s := []byte("")
	if config != "unset" {
		s = []byte(config)
	}
	err := os.WriteFile(fmt.Sprintf("%s/.kxd", loc), s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigs(configFileLocation string, homeDir string) []string {
	var files []string
	fileExt := getenv("KXD_MATCHER", ".conf")
	err := filepath.Walk(configFileLocation, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() && strings.Contains(f.Name(), fileExt) {
			files = append(files, f.Name())
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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func checkError(err error) {
	if err.Error() == "^D" {
		// https://github.com/manifoldco/promptui/issues/179
		log.Fatalf("<Del> not supported")
	} else if err.Error() == "^C" {
		os.Exit(1)
	} else {
		log.Fatal(err)
	}
}

type bellSkipper struct{}

func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
