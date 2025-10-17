package display

import (
    "fmt"
    "kmget/pkg/client"
    "kmget/pkg/configmap"
)

// PrintClusterInfo displays cluster information in a formatted way
func PrintClusterInfo(info *client.ClusterInfo) {
    fmt.Println("═══════════════════════════════════════════════════════════")
    fmt.Printf("Connected to Kubernetes Cluster\n")
    fmt.Println("═══════════════════════════════════════════════════════════")
    fmt.Printf("  Context:    %s\n", info.Context)
    fmt.Printf("  Cluster:    %s\n", info.Cluster)
    fmt.Printf("  Endpoint:   %s\n", info.Endpoint)
    fmt.Printf("  Namespace:  %s\n", info.Namespace)
    fmt.Printf("  Version:    %s\n", info.Version)
    fmt.Println("═══════════════════════════════════════════════════════════")
    fmt.Println()
}

// PrintConfigMapsList displays a list of ConfigMaps
func PrintConfigMapsList(namespace string, configMaps []configmap.ConfigMapInfo) {
    fmt.Printf("ConfigMaps in namespace '%s':\n", namespace)
    for _, cm := range configMaps {
        fmt.Printf("  - %s (data: %d, binary: %d)\n", cm.Name, cm.DataCount, cm.BinaryCount)
        for _, key := range cm.DataKeys {
            fmt.Printf("    * %s (text)\n", key)
        }
        for _, key := range cm.BinaryKeys {
            fmt.Printf("    * %s (binary)\n", key)
        }
    }
}

// PrintAllConfigMapsList displays ConfigMaps from all namespaces
func PrintAllConfigMapsList(allConfigMaps map[string][]configmap.ConfigMapInfo) {
    fmt.Println("ConfigMaps across all namespaces:")
    for namespace, configMaps := range allConfigMaps {
        if len(configMaps) > 0 {
            fmt.Printf("\nNamespace: %s\n", namespace)
            for _, cm := range configMaps {
                fmt.Printf("  - %s (data: %d, binary: %d)\n", cm.Name, cm.DataCount, cm.BinaryCount)
                for _, key := range cm.DataKeys {
                    fmt.Printf("    * %s (text)\n", key)
                }
                for _, key := range cm.BinaryKeys {
                    fmt.Printf("    * %s (binary)\n", key)
                }
            }
        }
    }
}

// PrintPullResult displays the result of pulling a ConfigMap
func PrintPullResult(result *configmap.PullConfigMapResult) {
    fmt.Printf("Pulling ConfigMap '%s' from namespace '%s':\n", result.ConfigMapName, result.Namespace)
    
    successCount := 0
    for _, file := range result.SavedFiles {
        if file.Success {
            if file.Binary {
                fmt.Printf("  ✓ Saved (binary): %s\n", file.Path)
            } else {
                fmt.Printf("  ✓ Saved: %s\n", file.Path)
            }
            successCount++
        } else {
            fmt.Printf("  ✗ Failed to save: %s (error: %v)\n", file.Path, file.Error)
        }
    }
    
    fmt.Printf("\nSuccessfully pulled %d/%d configuration file(s)\n", successCount, result.TotalFiles)
}

// PrintPullAllResults displays the results of pulling all ConfigMaps
func PrintPullAllResults(results []configmap.PullConfigMapResult) {
    totalConfigMaps := len(results)
    totalFiles := 0
    successfulFiles := 0

    fmt.Printf("Found %d ConfigMap(s) to process\n", totalConfigMaps)
    fmt.Println("Pulling ConfigMaps from all namespaces:")

    for _, result := range results {
        fmt.Printf("\n[Namespace: %s] ConfigMap: %s (%d files)\n", 
            result.Namespace, result.ConfigMapName, result.TotalFiles)
        
        for _, file := range result.SavedFiles {
            totalFiles++
            if file.Success {
                if file.Binary {
                    fmt.Printf("  ✓ Saved (binary): %s\n", file.Path)
                } else {
                    fmt.Printf("  ✓ Saved: %s\n", file.Path)
                }
                successfulFiles++
            } else {
                fmt.Printf("  ✗ Failed to save: %s (error: %v)\n", file.Path, file.Error)
            }
        }
    }

    fmt.Printf("\nSummary:\n")
    fmt.Printf("  - Processed %d ConfigMap(s)\n", totalConfigMaps)
    fmt.Printf("  - Successfully saved %d/%d configuration file(s)\n", successfulFiles, totalFiles)
}