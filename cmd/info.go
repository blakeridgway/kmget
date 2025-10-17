package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "kmget/pkg/client"
    "kmget/pkg/display"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
    Use:   "info",
    Short: "Display cluster information",
    Long:  `Display detailed information about the current Kubernetes cluster connection.`,
    Run: func(cmd *cobra.Command, args []string) {
        k8sClient, err := client.NewClient(kubeconfig)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
            os.Exit(1)
        }

        info, err := k8sClient.GetClusterInfo(kubeconfig)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error retrieving cluster info: %v\n", err)
            os.Exit(1)
        }

        display.PrintClusterInfo(info)
    },
}

func init() {
    rootCmd.AddCommand(infoCmd)
}