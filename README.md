# ConfiguRJ
ConfiguRJ is a Kubernetes Operator, that creates and updates Secrets and ConfigMaps in the cluster.
It is used, to copy one or more Secrets/ConfigMaps from one Namespace to another Namespace and keep 
the resources updated. 

## Table of Contents

- [Installation](#installation)
  - [Deploying to Kubernetes](#deploying-to-kubernetes)
  - [Configuration](#configuration)
- [Secrets and ConfigMaps](#secrets-and-configmaps)
  - [Original Annotations](#original-annotations)
  - [Replica Annotations](#replica-annotations)
- [Processes](#communication)
  - [Namespace Collection](#namespace-collection)
  - [Original Handling](#original-handling)
  - [Replica Distributing](#replica-distributing)
  - [Health-Operator](#health-operator)
  
## Installation
  
This part is about the installation of the ConfiguRJ service. A collection of the kubernetes manifests and 
a short explanation about the overall service configuration.
  
### Deploying to Kubernetes
  
To deploy the service to your cluster, there are the following manifests, which are recommended to run the service.
The manifests are minimalistic and do only contain the minimum neccessary information:
- Namespace
- ServiceAccount
- ClusterRole
- ClusterRoleBinding
- ConfigMap
- Deployment
  
#### Namespace
```
apiVersion: v1
kind: Namespace
metadata:
  name: configurj
---
```  
#### ServiceAccount
```
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: configurj
  name: configurj-sa
---
```  
#### ClusterRole
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cr-configurj
rules:
  # Get/List Namespaces
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["list", "get"]
  # Get/Create/List/Delete Conigmaps and Secrets
  - apiGroups: [""]
    resources: ["configmaps", "secrets"]
    verbs: ["list", "get", "create", "delete"]
---
```  
#### ClusterRoleBinding
```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: configurj-crb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cr-configurj
subjects:
  - kind: ServiceAccount
    name: configurj-sa
    namespace: configurj
---
```  
#### ConfigMap
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: configurj-settings
  namespace: configurj
data:
  settings.json: |
    "..."
---
```  
#### Deployment  
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: configurj
  namespace: configurj
  labels:
    app: configurj
spec:
  selector:
    matchLabels:
      app: configurj
  template:
    metadata:
      labels:
        app: configurj
    spec:
      serviceAccountName: configurj-sa
      containers:
      - name: configurj-controller
        image: docker.io/jnnkrdb/configurj:v1-stable
        imagePullPolicy: Always
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
          requests:
            memory: "64Mi"
            cpu: "250m"
        livenessProbe:
          initialDelaySeconds: 5
          periodSeconds: 5
          httpGet:
            path: /livez
            # the port will be configured in the configmap -> settings.json
            port: 8080
          failureThreshold: 10
        volumeMounts:
          - name: settings
            mountPath: /configs
            readOnly: true
      volumes:
        - name: settings
          configMap:
            name: configurj-settings
---
```  

### Configuration

This is the necessary content for the settings.json. 

```
{
    "immutablereplicas":true,
    "liveness":"8080",
    "sourcenamespace":"global-configs",
    "avoidsecrets":[
    "argocd",
    "configurj",
    "default",
    "elastic-system",
    "global-resources",
    "grafana",
    "infra",
    "kube-node-lease",
    "kube-public",
    "kube-system",
    "kubernetes-dashboard",
    "prometheus",
    "storage-nfs-ceph",
    "storage-rbd-ceph",
    "vaultwarden"
    ],
    "avoidconfigmaps":[
    "argocd",
    "configurj",
    "default",
    "elastic-system",
    "global-resources",
    "grafana",
    "infra",
    "kube-node-lease",
    "kube-public",
    "kube-system",
    "kubernetes-dashboard",
    "prometheus",
    "storage-nfs-ceph",
    "storage-rbd-ceph",
    "vaultwarden"
    ]
}

```

## Secrets and ConfigMaps
### Original Annotations
### Replica Annotations
  
## Processes
### Namespace Collection
### Original Handling
### Replica Distributing
### Health-Operator