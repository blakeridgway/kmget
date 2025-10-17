package client

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config holds the Kubernetes client configuration
type Config struct {
    Kubeconfig    string
    Namespace     string
    ConfigMap     string
    OutputDir     string
    ListOnly      bool
    AllNamespaces bool
}

// Client wraps the Kubernetes clientset with additional functionality
type Client struct {
    Clientset *kubernetes.Clientset
    Config    *Config
}

// NewClient creates a new Kubernetes client
func NewClient(kubeconfig string) (*Client, error) {
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, fmt.Errorf("failed to build config: %w", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create clientset: %w", err)
    }

    return &Client{
        Clientset: clientset,
    }, nil
}

// GetDefaultKubeconfig returns the default kubeconfig path
func GetDefaultKubeconfig() string {
    if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
        return kubeconfigEnv
    }
    if home := homedir.HomeDir(); home != "" {
        return filepath.Join(home, ".kube", "config")
    }
    return ""
}

// ClusterInfo represents cluster information
type ClusterInfo struct {
    Context     string
    Cluster     string
    Endpoint    string
    Namespace   string
    Version     string
    KubeconfigPath string
}

// GetClusterInfo retrieves cluster information
func (c *Client) GetClusterInfo(kubeconfigPath string) (*ClusterInfo, error) {
    version, err := c.Clientset.Discovery().ServerVersion()
    if err != nil {
        return nil, fmt.Errorf("failed to get server version: %w", err)
    }

    loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
    configOverrides := &clientcmd.ConfigOverrides{}
    kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

    rawConfig, err := kubeConfig.RawConfig()
    if err != nil {
        return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
    }

    currentContext := rawConfig.CurrentContext
    context, exists := rawConfig.Contexts[currentContext]
    var clusterName, namespace string
    if exists {
        clusterName = context.Cluster
        namespace = context.Namespace
        if namespace == "" {
            namespace = "default"
        }
    }

    cluster, exists := rawConfig.Clusters[clusterName]
    var endpoint string
    if exists {
        endpoint = cluster.Server
    }

    return &ClusterInfo{
        Context:        currentContext,
        Cluster:        clusterName,
        Endpoint:       endpoint,
        Namespace:      namespace,
        Version:        version.GitVersion,
        KubeconfigPath: kubeconfigPath,
    }, nil
}

// GetNamespaces returns all namespaces
func (c *Client) GetNamespaces() ([]string, error) {
    ctx := context.Background()
    namespaces, err := c.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list namespaces: %w", err)
    }

    var names []string
    for _, ns := range namespaces.Items {
        names = append(names, ns.Name)
    }
    return names, nil
}