# ConfiguRJ
ConfiguRJ is a Kubernetes Operator Package, that creates and updates Secrets and ConfigMaps in the cluster. It is used, to read a GlobalConfig or a GlobalSecret from the cluster and replicate its data into ConfigMaps or Secrets in the desired namespaces and keep the resources updated. 

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/jnnkrdb/configurj)
[![CodeFactor](https://www.codefactor.io/repository/github/jnnkrdb/configurj/badge)](https://www.codefactor.io/repository/github/jnnkrdb/configurj)
[![Go Report Card](https://goreportcard.com/badge/github.com/jnnkrdb/configurj)](https://goreportcard.com/report/github.com/jnnkrdb/configurj)
[![GitHub issues](https://badgen.net/github/issues/jnnkrdb/configurj/)](https://github.com/jnnkrdb/configurj/issues/)

## Table of Contents

- [Installation](#installation)
  - [Defaults](#defaults)
    - [Namespace](#namespace)
    - [CustomResourceDefinition GlobalConfig](#customresourcedefinition-globalconfig)
    - [CustomResourceDefinition GlobalSecret](#customresourcedefinition-globalsecret)
  - [Operator](#operator)
    - [ServiceAccount](#serviceaccount)
    - [ClusterRole](#clusterrole)
    - [ClusterRoleBinding](#clusterrolebinding)
    - [Deployment](#deployment)
    - [CustomResourceDefinition](#customresourcedefinition)
  - [Example Deployments](#example-deployments)
    - [GlobalConfig](#globalconfig)
    - [GlobalSecret](#globalsecret)
- [Configuration](#configuration)
  - [Operator Environment Variables](#operator-environment-variables)
  - [UI-Controller Angular Config](#ui-controller-angular-config)
- [RoadMap or Planned](#roadmap-or-planned)
  
## Installation
  
This part is about the installation of the ConfiguRJ service. It contains the collection of the kubernetes manifests and a short explanation about the overall service configuration. To get this service running, you need to deploy the yaml-files to your kubernetes cluster. The deployment of the ConfigMaps/Secrets will be handled with the CRDs of this project. Deploy a GlobalConfig to rollout ConfigMaps into the configured namespaces or deploy a GlobalSecret to do so with Secrets.
  
To deploy the service to your cluster, there are the following manifests, which are recommended to run the service. The manifests are minimalistic and do only contain the minimum neccessary information.

### Defaults
All ConfiguRJ-Controllers need some default deployment-manifests, for example the CustomResourceDefinitions. Those default deployments are listed in this section.
- [Namespace](#namespace)
- [CustomResourceDefinition GlobalConfig](#customresourcedefinition-globalconfig)
- [CustomResourceDefinition GlobalSecret](#customresourcedefinition-globalsecret)

#### Namespace
```yaml
---
apiVersion: v1
kind: Namespace
metadata:
  name: configurj
  labels:
    app: configurj
```  
#### CustomResourceDefinition GlobalConfig
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalconfigs.globals.jnnkrdb.de
  labels:
    app: configurj
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
                  default: false
                namespaces: 
                  type: object
                  properties:
                    avoidregex:
                      type: array
                      items: 
                        type: string
                      default: []
                    matchregex:
                      type: array
                      items: 
                        type: string
                      default: []
                  required: [avoidregex, matchregex]
                name: 
                  type: string
                data:
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
                  default: {}
              required: [name, namespaces]
```  
#### CustomResourceDefinition GlobalSecret
```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalsecrets.globals.jnnkrdb.de
  labels:
    app: configurj
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
              # crds properties
              properties:
                immutable: 
                  type: boolean
                  default: false
                namespaces: 
                  type: object
                  properties:
                    avoidregex:
                      type: array
                      items: 
                        type: string
                      default: []
                    matchregex:
                      type: array
                      items: 
                        type: string
                      default: []
                  required: [avoidregex, matchregex]
                name: 
                  type: string
                type: 
                  type: string
                  enum:
                    - Opaque
                    - kubernetes.io/service-account-token
                    - kubernetes.io/dockercfg
                    - kubernetes.io/dockerconfigjson
                    - kubernetes.io/basic-auth
                    - kubernetes.io/ssh-auth
                    - kubernetes.io/tls
                    - bootstrap.kubernetes.io/token
                data:
                  type: object
                  x-kubernetes-preserve-unknown-fields: true
                  default: {}
              required: [name, namespaces, type]
```

### Operator
The Operator contains the core functionality of this controller package. The operator requests GlobalConfigs and GlobalSecrets in the cluster and creates the ConfigMaps and Secrets with their specifications.
The Controller needs some specific kubernetes manifests to show its full potential:
  - [ServiceAccount](#serviceaccount)
  - [ClusterRole](#clusterrole)
  - [ClusterRoleBinding](#clusterrolebinding)
  - [Deployment](#deployment)

#### ServiceAccount
```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: configurj
  name: sa-configurj-operator
  labels:
    app: configurj
    type: operator
```  
#### ClusterRole
```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cr-configurj-operator
  labels:
    app: configurj
    type: operator
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
  name: crb-configurj-operator
  labels:
    app: configurj
    type: operator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cr-configurj-operator
subjects:
  - kind: ServiceAccount
    name: sa-configurj-operator
    namespace: configurj
```  
#### Deployment
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: configurj-operator
  namespace: configurj
  labels:
    app: configurj
    type: operator
spec:
  selector:
    matchLabels:
      app: configurj
      type: operator
  template:
    metadata:
      labels:
        app: configurj
        type: operator
    spec:
      serviceAccountName: sa-configurj-operator
      containers:
      - name: configurj-operator
        image: docker.io/jnnkrdb/configurj-engine:latest
        imagePullPolicy: Always
        env:
          - name: DEBUGGING
            value: "true"
          - name: TIMEOUTMINUTES
            value: "3"
        resources:
          limits:
            memory: "128Mi"
            cpu: "500m"
          requests:
            memory: "64Mi"
            cpu: "250m"
        livenessProbe:
          initialDelaySeconds: 10
          periodSeconds: 5
          httpGet:
            path: /healthz/live
            port: 80
          failureThreshold: 5
```

### Example Deployments

In this section you can find some example deployments of the GlobalConfig and/or GlobalSecret resources.
  - [GlobalConfig](#globalconfig)
  - [GlobalSecret](#globalsecret)

#### GlobalConfig
```yaml
---
apiVersion: globals.jnnkrdb.de/v1alpha1
kind: GlobalConfig
metadata:
  name: gc-name
  namespace: default
spec:
  immutable: false
  namespaces:
    avoidregex: 
      - default # matches namespace "default" -> namespace default will be avoided
      - prod. # matches namespaces like "production-financial", "prod-databases", "prod*" -> namespaces like "production-financial", "prod-databases" or "prod*" will be avoided
    matchregex: 
      - production-mssql # matches namespace "production-mssql", BUT since "prod."-regex is in the avoidregex-list, this namespace will not be matched
      - .dev # matches namespaces like "financials-dev", "databases-dev", "dev", etc. -> namespaces with the suffix "dev" will be matched
      - .internal. # matches namespaces like "test-internal-financials", "databases-internals", "internal", etc. -> namespaces, which contain the substring "internal" will be matched
  name: cm-name
  data: # the data section should be filled like the data-section of a normal configmap

    # kubernetes example of a configmap -> https://kubernetes.io/docs/concepts/configuration/configmap/
    # property-like keys; each key maps to a simple value
    player_initial_lives: "3"
    ui_properties_file_name: "user-interface.properties"

    # file-like keys
    game.properties: |
      enemy.types=aliens,monsters
      player.maximum-lives=5    
    user-interface.properties: |
      color.good=purple
      color.bad=yellow
      allow.textmode=true    

```  
#### GlobalSecret
```yaml
---
apiVersion: globals.jnnkrdb.de/v1alpha1
kind: GlobalSecret
metadata:
  name: gs-name
spec:
  immutable: true
  namespaces:
    avoidregex: []
    matchregex: 
      - "." # matches all namespaces
  name: scrt-name
  type: kubernetes.io/dockerconfigjson # or other type, supported by kubernetes secrets -> https://kubernetes.io/docs/concepts/configuration/secret/
  data: # must be base64 encrypted by yourself, but like the globalconfig, this section is build like its underlying secret
    .dockerconfigjson: <base64 encrypted docker config json file>
```

## Configuration

The Operator package must be configured for each controller seperatly.
  - [Operator Environment Variables](#operator-environment-variables)
  - [UI-Controller Angular Config](#ui-controller-angular-config)

#### Operator Environment Variables

- `DEBUGGING` (+Optional): This env-variable configures the print output to the controller-commandline. If set to `"true"`, more detailed objects will be printed and not only the error-messages. The default is: `false`
- `TIMEOUTMINUTES` (+Optional): This variable sets the await timeout of the routine in minutes. Its soure code counterpart is an float64 type. The default is: `3` minutes.

#### UI-Controller Angular Config

Under construction...

## RoadMap or Planned
- Ingress Configuration
- Angular UI -> Overview of Globals
  - ``docker.io/jnnkrdb/configurj-ui:latest``
- BackendAPI-Controller feeding Angular UI
- Better integrated Health probes
