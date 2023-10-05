package utils

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"log"
	"sort"
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
