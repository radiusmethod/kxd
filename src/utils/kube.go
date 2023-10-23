package utils

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func InitializeKubeconfig(kubeconfigPath string) (*api.Config, error) {
	return clientcmd.LoadFromFile(kubeconfigPath)
}

func SwitchContext(config *api.Config, contextName string, kubeconfigPath string) error {
	config.CurrentContext = contextName
	return clientcmd.WriteToFile(*config, kubeconfigPath)
}

func ListContexts(kubeconfigPath string) []string {
	config, err := InitializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	var contexts []string
	for contextName := range config.Contexts {
		contexts = append(contexts, contextName)
	}
	sort.Strings(contexts)
	return contexts
}

func GetConfigs() []string {
	var files []string
	configFileLocation := GetConfigFileLocation()

	fileExts := strings.Split(GetEnv("KXD_MATCHER", ".conf"), ",")
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

	defaultConfigPath := filepath.Join(GetHomeDir(), ".kube/config")
	if _, err := os.Stat(defaultConfigPath); err == nil {
		files = append(files, "default")
	}
	files = append(files, "unset")
	sort.Strings(files)
	return files
}
