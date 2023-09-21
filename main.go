package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	NoticeColor = "\033[0;38m%s\u001B[0m"
	PromptColor = "\033[1;38m%s\u001B[0m"
	CyanColor   = "\033[0;36m%s\033[0m"
)

var version string = "v0.0.1"

func newPromptUISearcher(items []string) list.Searcher {
	return func(searchInput string, itemIndex int) bool {
		return strings.Contains(strings.ToLower(items[itemIndex]), strings.ToLower(searchInput))
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("kxd version", version)
		os.Exit(0)
	}
	home := os.Getenv("HOME")
	configFileLocation := fmt.Sprintf("%s/.kube", home)
	configs := getConfigs(configFileLocation)

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
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf(PromptColor, "Choose a config")
	fmt.Printf(NoticeColor, "? ")
	fmt.Printf(CyanColor, result)
	fmt.Println("")
	writeFile(result, home)
}

func writeFile(profile, loc string) {
	s := []byte("")
	if profile != "default" {
		s = []byte(profile)
	}
	err := os.WriteFile(fmt.Sprintf("%s/.kxd", loc), s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigs(configFileLocation string) []string {
	var files []string
	filepath.Walk(configFileLocation, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".conf" ||
				(os.Getenv("KXD_MATCHER") != "" && strings.Contains(f.Name(), os.Getenv("KXD_MATCHER"))) {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	files = append(files, "default")
	sort.Strings(files)
	return files
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
