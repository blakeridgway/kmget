package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "kmget/pkg/client"
)

var (
    cfgFile       string
    kubeconfig    string
    namespace     string
    outputDir     string
    allNamespaces bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "kmget",
    Short: "A CLI tool to pull configuration files from Kubernetes ConfigMaps",
    Long: `ConfigMap Puller is a CLI tool that helps you retrieve configuration files
from Kubernetes ConfigMaps. It supports listing ConfigMaps and pulling their
contents to local files.

Examples:
  # List ConfigMaps in default namespace
  kmget list

  # List all ConfigMaps across all namespaces
  kmget list --all-namespaces

  # Pull a specific ConfigMap
  kmget pull my-config --namespace default

  # Pull all ConfigMaps from all namespaces
  kmget pull --all-namespaces`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kmget.yaml)")
    rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", client.GetDefaultKubeconfig(), "path to kubeconfig file (respects KUBECONFIG env var)")
    rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
    rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "output directory for config files")
    rootCmd.PersistentFlags().BoolVar(&allNamespaces, "all-namespaces", false, "operate on all namespaces")

    viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
    viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
    viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
    viper.BindPFlag("all-namespaces", rootCmd.PersistentFlags().Lookup("all-namespaces"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)

        viper.AddConfigPath(home)
        viper.SetConfigType("yaml")
        viper.SetConfigName(".kmget")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}