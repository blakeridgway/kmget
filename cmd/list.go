package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "kmget/pkg/client"
    "kmget/pkg/configmap"
    "kmget/pkg/display"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List ConfigMaps",
    Long: `List ConfigMaps in the specified namespace or across all namespaces.

Examples:
  # List ConfigMaps in default namespace
  kmget list

  # List ConfigMaps in specific namespace
  kmget list --namespace kube-system

  # List ConfigMaps across all namespaces
  kmget list --all-namespaces`,
    Run: func(cmd *cobra.Command, args []string) {
        k8sClient, err := client.NewClient(kubeconfig)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
            os.Exit(1)
        }

        ops := configmap.NewOperations(k8sClient.Clientset)

        if allNamespaces {
            allConfigMaps, err := ops.ListAllConfigMaps()
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error listing ConfigMaps: %v\n", err)
                os.Exit(1)
            }
            display.PrintAllConfigMapsList(allConfigMaps)
        } else {
            configMaps, err := ops.ListConfigMaps(namespace)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error listing ConfigMaps: %v\n", err)
                os.Exit(1)
            }
            display.PrintConfigMapsList(namespace, configMaps)
        }
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}