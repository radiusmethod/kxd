package utils

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func ListNamespaces(kubeconfigPath string) []string {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("error creating Kubernetes client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var namespaces []string
	nss, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("error listing namespaces: %w", err)
	}

	for _, ns := range nss.Items {
		namespaces = append(namespaces, ns.Name)
	}

	sort.Strings(namespaces)
	return namespaces
}

func SwitchNamespace(config *api.Config, namespaceName string, kubeconfigPath string) error {
	config, err := InitializeKubeconfig(kubeconfigPath)
	if err != nil {
		log.Fatalf("Error initializing kubeconfig: %v\n", err)
	}

	contextName := config.CurrentContext
	context, exists := config.Contexts[contextName]
	if !exists {
		log.Fatalf("Context %s does not exist in kubeconfig", contextName)
	}
	context.Namespace = namespaceName

	return clientcmd.WriteToFile(*config, kubeconfigPath)
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
