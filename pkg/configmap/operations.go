package configmap

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// Operations handles ConfigMap operations
type Operations struct {
    clientset *kubernetes.Clientset
}

// NewOperations creates a new ConfigMap operations handler
func NewOperations(clientset *kubernetes.Clientset) *Operations {
    return &Operations{
        clientset: clientset,
    }
}

// ConfigMapInfo represents ConfigMap information
type ConfigMapInfo struct {
    Name        string
    Namespace   string
    DataKeys    []string
    BinaryKeys  []string
    DataCount   int
    BinaryCount int
}

// GetConfigMap retrieves a specific ConfigMap
func (o *Operations) GetConfigMap(namespace, name string) (*corev1.ConfigMap, error) {
    ctx := context.Background()
    configMap, err := o.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get ConfigMap '%s' in namespace '%s': %w", name, namespace, err)
    }
    return configMap, nil
}

// ListConfigMaps lists all ConfigMaps in a namespace
func (o *Operations) ListConfigMaps(namespace string) ([]ConfigMapInfo, error) {
    ctx := context.Background()
    configMaps, err := o.clientset.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list ConfigMaps in namespace '%s': %w", namespace, err)
    }

    var infos []ConfigMapInfo
    for _, cm := range configMaps.Items {
        var dataKeys, binaryKeys []string
        for key := range cm.Data {
            dataKeys = append(dataKeys, key)
        }
        for key := range cm.BinaryData {
            binaryKeys = append(binaryKeys, key)
        }

        infos = append(infos, ConfigMapInfo{
            Name:        cm.Name,
            Namespace:   cm.Namespace,
            DataKeys:    dataKeys,
            BinaryKeys:  binaryKeys,
            DataCount:   len(cm.Data),
            BinaryCount: len(cm.BinaryData),
        })
    }

    return infos, nil
}

// ListAllConfigMaps lists ConfigMaps from all namespaces
func (o *Operations) ListAllConfigMaps() (map[string][]ConfigMapInfo, error) {
    ctx := context.Background()
    namespaces, err := o.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to list namespaces: %w", err)
    }

    result := make(map[string][]ConfigMapInfo)
    for _, ns := range namespaces.Items {
        configMaps, err := o.ListConfigMaps(ns.Name)
        if err != nil {
            return nil, err
        }
        if len(configMaps) > 0 {
            result[ns.Name] = configMaps
        }
    }

    return result, nil
}

// SaveResult represents the result of saving a file
type SaveResult struct {
    Path    string
    Success bool
    Error   error
    Binary  bool
}

// PullConfigMapResult represents the result of pulling a ConfigMap
type PullConfigMapResult struct {
    ConfigMapName string
    Namespace     string
    SavedFiles    []SaveResult
    TotalFiles    int
}

// PullConfigMap saves a ConfigMap's data to files
func (o *Operations) PullConfigMap(namespace, name, outputDir string) (*PullConfigMapResult, error) {
    configMap, err := o.GetConfigMap(namespace, name)
    if err != nil {
        return nil, err
    }

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create output directory: %w", err)
    }

    result := &PullConfigMapResult{
        ConfigMapName: name,
        Namespace:     namespace,
        SavedFiles:    []SaveResult{},
    }

    // Handle text data
    for key, value := range configMap.Data {
        outputPath := filepath.Join(outputDir, key)
        saveResult := SaveResult{
            Path:   outputPath,
            Binary: false,
        }

        if err := os.WriteFile(outputPath, []byte(value), 0644); err != nil {
            saveResult.Success = false
            saveResult.Error = err
        } else {
            saveResult.Success = true
        }

        result.SavedFiles = append(result.SavedFiles, saveResult)
        result.TotalFiles++
    }

    // Handle binary data
    for key, value := range configMap.BinaryData {
        outputPath := filepath.Join(outputDir, key)
        saveResult := SaveResult{
            Path:   outputPath,
            Binary: true,
        }

        if err := os.WriteFile(outputPath, value, 0644); err != nil {
            saveResult.Success = false
            saveResult.Error = err
        } else {
            saveResult.Success = true
        }

        result.SavedFiles = append(result.SavedFiles, saveResult)
        result.TotalFiles++
    }

    return result, nil
}

// PullAllConfigMaps saves all ConfigMaps from all namespaces
func (o *Operations) PullAllConfigMaps(outputDir string) ([]PullConfigMapResult, error) {
    allConfigMaps, err := o.ListAllConfigMaps()
    if err != nil {
        return nil, err
    }

    var results []PullConfigMapResult
    for namespace, configMaps := range allConfigMaps {
        for _, cm := range configMaps {
            if cm.DataCount == 0 && cm.BinaryCount == 0 {
                continue // Skip empty ConfigMaps
            }

            nsDir := filepath.Join(outputDir, namespace)
            result, err := o.PullConfigMap(namespace, cm.Name, nsDir)
            if err != nil {
                return results, fmt.Errorf("failed to pull ConfigMap '%s' from namespace '%s': %w", cm.Name, namespace, err)
            }
            results = append(results, *result)
        }
    }

    return results, nil
}