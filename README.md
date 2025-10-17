# kmget

A command-line tool for retrieving configuration files from Kubernetes ConfigMaps in AKS (Azure Kubernetes Service) and other Kubernetes clusters.

## Features

- 🔍 List ConfigMaps in namespaces or across all namespaces
- 📁 Pull ConfigMap data to local files
- 🌐 Multi-namespace support
- 📊 Display cluster connection information
- 🔄 Handle both text and binary ConfigMap data
- ⚙️ Flexible kubeconfig and in-cluster authentication

## Installation

### Build from Source

```bash
git clone https://github.com/blakeridgway/kmget
cd kmget
go build -o kmget
```

### Prerequisites

- Go 1.21+
- Kubernetes cluster access
- Valid kubeconfig file

## Quick Start

```bash
# Display cluster info
kmget info

# List ConfigMaps
kmget list --all-namespaces

# Pull a specific ConfigMap
kmget pull my-config -n default -o ./config

# Pull all ConfigMaps
kmget pull --all-namespaces -o ./all-configs
```

## Commands

### `kmget info`
Display cluster connection information.

### `kmget list`
List ConfigMaps in namespace or all namespaces.

```bash
kmget list -n kube-system
kmget list --all-namespaces
```

### `kmget pull [CONFIGMAP_NAME]`
Pull ConfigMap data to local files.

```bash
kmget pull app-config -n default -o ./config
kmget pull --all-namespaces -o ./backup
```

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--kubeconfig` | | `~/.kube/config` | Path to kubeconfig file |
| `--namespace` | `-n` | `default` | Kubernetes namespace |
| `--output` | `-o` | `.` | Output directory |
| `--all-namespaces` | | `false` | Operate on all namespaces |

## Examples

### Azure Kubernetes Service (AKS)

```bash
# Connect to AKS
az aks get-credentials --resource-group myRG --name myCluster

# Backup all configurations
kmget pull --all-namespaces --output ./backup
```

### Multiple Environments

```bash
# Development
export KUBECONFIG=~/.kube/dev-config
kmget pull --all-namespaces -o ./dev-configs

# Production
export KUBECONFIG=~/.kube/prod-config
kmget pull --all-namespaces -o ./prod-configs
```

## Configuration

### Environment Variables
- `KUBECONFIG`: Path to kubeconfig file

### Config File
Create `~/.kmget.yaml`:
```yaml
kubeconfig: /path/to/kubeconfig
namespace: default
output: ./output
```

## Output Structure

```
output-directory/
├── namespace1/
│   ├── configmap1_key1.yaml
│   └── configmap1_key2.json
└── namespace2/
    └── configmap2_config.properties
```

## Troubleshooting

**Authentication Issues:**
```bash
kubectl config current-context
kmget info
```

**Permission Issues:**
```bash
kubectl auth can-i get configmaps
```

**RBAC Requirements:**
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kmget-reader
rules:
- apiGroups: [""]
  resources: ["configmaps", "namespaces"]
  verbs: ["get", "list"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.