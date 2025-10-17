package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "kmget/pkg/client"
    "kmget/pkg/configmap"
    "kmget/pkg/display"
)

var (
    configMapName string
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
    Use:   "pull [CONFIGMAP_NAME]",
    Short: "Pull ConfigMap data to local files",
    Long: `Pull ConfigMap data to local files in the specified output directory.

Examples:
  # Pull a specific ConfigMap
  kmget pull my-config --namespace default --output ./config

  # Pull all ConfigMaps from all namespaces
  kmget pull --all-namespaces --output ./all-configs

  # Pull ConfigMap using positional argument
  kmget pull my-config`,
    Args: func(cmd *cobra.Command, args []string) error {
        if !allNamespaces && len(args) == 0 && configMapName == "" {
            return fmt.Errorf("ConfigMap name is required when not using --all-namespaces flag")
        }
        return nil
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Use positional argument if provided
        if len(args) > 0 {
            configMapName = args[0]
        }

        k8sClient, err := client.NewClient(kubeconfig)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
            os.Exit(1)
        }

        ops := configmap.NewOperations(k8sClient.Clientset)

        if allNamespaces {
            results, err := ops.PullAllConfigMaps(outputDir)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error pulling ConfigMaps: %v\n", err)
                os.Exit(1)
            }
            display.PrintPullAllResults(results)
        } else {
            result, err := ops.PullConfigMap(namespace, configMapName, outputDir)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error pulling ConfigMap: %v\n", err)
                os.Exit(1)
            }
            display.PrintPullResult(result)
        }
    },
}

func init() {
    pullCmd.Flags().StringVarP(&configMapName, "configmap", "c", "", "name of the ConfigMap to pull")
    rootCmd.AddCommand(pullCmd)
}