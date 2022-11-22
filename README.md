# ConfiguRJ
ConfiguRJ is a Kubernetes Operator, that creates and updates Secrets and ConfigMaps in the cluster. It is used, to read a GlobalConfig or a GlobalSecret from the cluster and replicate its data intto ConfigMaps or Secrets in the desired namespaces and keep the resources updated. 

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/jnnkrdb/configurj)
[![CodeFactor](https://www.codefactor.io/repository/github/jnnkrdb/configurj/badge)](https://www.codefactor.io/repository/github/jnnkrdb/configurj)
[![Go Report Card](https://goreportcard.com/badge/github.com/jnnkrdb/configurj)](https://goreportcard.com/report/github.com/jnnkrdb/configurj)
[![GitHub issues](https://badgen.net/github/issues/jnnkrdb/configurj/)](https://github.com/jnnkrdb/configurj/issues/)

## Table of Contents

- [Installation](#installation)
  - [Deploying to Kubernetes](#deploying-to-kubernetes)
    - [Namespace](#namespace)
    - [ServiceAccount](#serviceaccount)
    - [ClusterRole](#clusterrole)
    - [ClusterRoleBinding](#clusterrolebinding)
    - [ConfigMap](#configmap)
    - [Deployment](#deployment)
    - [CustomResourceDefinition](#customresourcedefinition)
  - [Configuration](#configuration)
- [RoadMap or Planned](#roadmap-or-planned)
  
## Installation
  
This part is about the installation of the ConfiguRJ service. It contains the collection of the kubernetes manifests and a short explanation about the overall service configuration. To get this service running, you need to deploy the yaml-files to your kubernetes cluster. The deployment of the ConfigMaps/Secrets wll be handled with the CRDs of this project. Deploy a GlobalConfig to rollout ConfigMaps into the configured namespaces or deploy a GlobalSecret to do so with Secrets.
  
### Deploying to Kubernetes
  
To deploy the service to your cluster, there are the following manifests, which are recommended to run the service.

The manifests are minimalistic and do only contain the minimum neccessary information:
- Namespace
- ServiceAccount
- ClusterRole
- ClusterRoleBinding
- ConfigMap
- Deployment
- CustomResourceDefinition
  
#### Namespace
```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: configurj
```  
#### ServiceAccount
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: configurj
  name: configurj-sa
```  
#### ClusterRole
```yaml
---
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
  # Get/Create/List/Delete Conigmaps and Secrets
  - apiGroups: ["globals.jnnkrdb.de"]
    resources: ["globalconfigs", "globalsecrets"]
    verbs: ["get", "list"]
```  
#### ClusterRoleBinding
```yaml
---
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
```  
#### ConfigMap
```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: configurj-settings
  namespace: configurj
data:
  # Find the settings content below in the configuration topic
  settings.json: |
    "..." 
```  
#### Deployment  
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: configurj-engine
  namespace: configurj
  labels:
    app: configurj
    part: backend
spec:
  selector:
    matchLabels:
      app: configurj
      part: backend
  template:
    metadata:
      labels:
        app: configurj
        part: backend
    spec:
      serviceAccountName: configurj-sa
      containers:
      - name: configurj-engine
        image: docker.io/jnnkrdb/configurj-engine:v1beta1
        imagePullPolicy: Always
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
          requests:
            memory: "64Mi"
            cpu: "250m"
            # currently the healthprobe is under maintenance
#        livenessProbe:
#          initialDelaySeconds: 5
#          periodSeconds: 5
#          httpGet:
#            path: /livez
#            port: 8080
#          failureThreshold: 10
        volumeMounts:
          - name: settings
            mountPath: /configs
            readOnly: true
      volumes:
        - name: settings
          configMap:
            name: configurj-settings
            items:
              - key: settings.json
                path: settings.json
```  
#### CustomResourceDefinition
GlobalConfig.jnnkrdb.de/v1beta1
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalconfigs.globals.jnnkrdb.de
spec:
  group: globals.jnnkrdb.de
  scope: Namespaced
  names:
    plural: globalconfigs
    singular: globalconfig
    kind: GlobalConfig
    shortNames: 
    - gc
    - gcs
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema: 
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              # crds properties
              properties:
                immutable: 
                  type: boolean
                namespaces: 
                  type: array
                  items: 
                    type: string
                name: 
                  type: string
                data:
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
              required: [namespaces, name, data]
```  
GlobalSecret.jnnkrdb.de/v1beta1
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalsecrets.globals.jnnkrdb.de
spec:
  group: globals.jnnkrdb.de
  scope: Namespaced
  names:
    plural: globalsecrets
    singular: globalsecret
    kind: GlobalSecret
    shortNames: 
    - gs
    - gss
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema: 
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              x-kubernetes-validation: 
                - rule: "self.type in ['Opaque','kubernetes.io/service-account-token','kubernetes.io/dockercfg','kubernetes.io/dockerconfigjson','kubernetes.io/basic-auth','kubernetes.io/ssh-auth','kubernetes.io/tls','bootstrap.kubernetes.io/token']"
                  message: "please use an existing secret.type -> https://kubernetes.io/docs/concepts/configuration/secret/#secret-types"
              # crds properties
              properties:
                immutable: 
                  type: boolean
                namespaces: 
                  type: array
                  items: 
                    type: string
                name: 
                  type: string
                type: 
                  type: string
                data:
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
              required: [namespaces, name, data, type]
```

### Configuration
This is the necessary content for the settings.json. `debugging` is a bool-value, which determines, if the debbuging-print to console will be activated. `timeoutsec` is a float64-value, which determines the wait timeout between each routine.`globalavoidnamespaces` is a string-value collection of the namespaces in the cluster, that will be avoided in any case. Namespaces that will be configured in these avoids are on the highest avoid-priority.

```json
{
  "debugging":true,
  "timeoutsec": 5,
  "globalavoidnamespaces":[
    "default",
    "kube-system"
  ]
}
```

## RoadMap or Planned
- Ingress Configuration
- Angular UI -> Overview of Globals
  - ``docker.io/jnnkrdb/configurj-ui:latest``
- Matching namespaces with regex