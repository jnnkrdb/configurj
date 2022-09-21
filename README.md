# ConfiguRJ
ConfiguRJ is a Kubernetes Operator, that creates and updates Secrets and ConfigMaps 
in the cluster. It is used, to copy one or more Secrets/ConfigMaps from one Namespace 
to another Namespace and keep the resources updated. 

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/jnnkrdb/configurj)
[![CodeFactor](https://www.codefactor.io/repository/github/jnnkrdb/configurj/badge)](https://www.codefactor.io/repository/github/jnnkrdb/configurj)
[![Go Report Card](https://goreportcard.com/badge/github.com/jnnkrdb/configurj)](https://goreportcard.com/report/github.com/jnnkrdb/configurj)

## Table of Contents

- [Installation](#installation)
  - [Deploying to Kubernetes](#deploying-to-kubernetes)
  - [Configuration](#configuration)
- [Secrets and ConfigMaps](#secrets-and-configmaps)
  - [Original Annotations](#original-annotations)
  - [Replica Annotations](#replica-annotations)
- [Routines](#routines)
  - [Namespace Collection](#namespace-collection)
  - [Original Handling](#original-handling)
  - [Replica Distribution](#replica-distribution)
  - [Health-Operator](#health-operator)
  
## Installation
  
This part is about the installation of the ConfiguRJ service. A collection of the 
kubernetes manifests and a short explanation about the overall service configuration. 
To get this service running, you need to deploy the yaml-files to your kubernetes cluster
and then create a `<namespace>`, where the secrets and configmaps, you want to distribute to the 
whole cluster, will be deployed/stored. In the service configmap, that is stated below, you
have to set the `"sourcenamespace":"<namespace>"` value, to the created namespace, and the service 
will constantly look for secrets and configmaps with the required annotations.
  
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
  settings.json: |
    "..."
```  
#### Deployment  
```yaml
---
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
```  

### Configuration

This is the necessary content for the settings.json. 

`immutablereplicas` is a bool-value, which determines, if the replicas will be immutable or not. `healthport` 
is a string-value, which determines the port for the readyness and liveness probes. `sourcenamespace` is a 
string-value, which determines the namespace, that the operator will collect the object (Secrets, ConfigMaps) 
from. `avoidsecrets/avoidconfigmaps` are string-value collections of the namespaces in the cluster, that will 
be avoided in any case. Namespaces that will be configured in these avoids are on the highest avoid-priority.

```json
{
    "immutablereplicas":true,
    "healthport":"8080",
    "sourcenamespace":"global-configs",
    "avoidsecrets":[
    "default",
    "global-resources",
    "kube-node-lease",
    "kube-public",
    "kube-system"
    ],
    "avoidconfigmaps":[
    "default",
    "global-resources",
    "kube-node-lease",
    "kube-public",
    "kube-system"
    ]
}
```

## Secrets and ConfigMaps

The Service gets the necessary information per secret/configmap from the secret/configmap itself. To
get the informations, the service uses some annotations. The original and the replica get different 
annotations, some are necessary for the service, some give the administrator information about the object.

### Original Annotations

```yaml
configurj.jnnkrdb.de/active: "true"/"false"
```
This annotation marks the object, to be replicated. If "true", the object will be replicated, if "false" the 
object will be removed from the other namespaces. If the annotation doesn't exist, the object will be ignored 
completly.

```yaml
configurj.jnnkrdb.de/avoid: "namespace-1;namespace-2"
```
This annotation is a collection of the namespaces, that the object should avoid additionally to the global avoids
from the settings.json. Seperate the namespaces with ";". Use "*", to avoid all existing namespaces.

```yaml
configurj.jnnkrdb.de/match: "namespace-3;namespace-4"
```
This annotation is a collection of the namespaces, that the object should match without the global avoids
from the settings.json. Seperate the namespaces with ";". Use "*", to match all existing namespaces.

### Replica Annotations

```yaml
configurj.jnnkrdb.de/replica: "true"
```
This annotation is set to "true" by default from the service at creation time. It is used as a marker, to declare 
an item as a replica. If the annotation is removed, the replicated item will be ignored and not be updated or deleted.

```yaml
configurj.jnnkrdb.de/timestamp: "YYYY/MM/DD"
```
This annotation is an information for the administrator, to see the last time, the item was changed or created.
The annotation is not necessary for the service to handle the objects.

```yaml
configurj.jnnkrdb.de/original: "<original-name>"
```
The name of the original resource is stored in this annotation. 

```yaml
configurj.jnnkrdb.de/original-ns: "<original-namespace>"
```
The namespace of the original resource is stored in this annotation. 

```yaml
configurj.jnnkrdb.de/original-rv: "<original-resourceversion>"
```
The resourceversion of the original resource is stored in this annotation. If the resourceversion in this annotation differs
from the resourceversion of the original item while comparision, the replica will be updated to the new version. 

## Routines

This service is build out of several routines. Two of them handle the replication of the entities ConfigMap and Secret and the
comparision of the replica and the original. One routine handles the collection of the namespaces, in which the entities will be
deployed, in a clusterwide consideration. The last routine handles the health probes, to provide kubernetes with the necessary 
health informations about the service.

### Namespace Collection

The namespace collection routine handles the repeated gathering of all namespaces in the cluster. This routine runs in an infinite
loop, to constantly get all current namespaces. If the routine failes, the health-state of the service will be set to unhealthy (Code 500) 
and the cached list of the namespaces will be replaced with an empty list. 


The collected namespaces will then
be compared with the configured, clusterwide avoidances, per objecttype (Secret, ConfigMap), to create namespace

### Original Handling

The "original handling" is one part of a routine, where the original of a specific object type, for example a Secret, will be determined
and the information about the object will be collected. The information collection contains the data about the object that will be replicated,
lists of the namespaces which should be avoided for this specific object, or which should be match and if the replication should be executed or 
skipped.

### Replica Distributing

The "replica distribution" is the second part of a routine, where the replicas of a specific object type will be created. The provided information
from the first part of the routine ("Original Handling") is used to create the replicas and deploy them to the configured namespaces. 

### Health-Operator

This process handles the health probes requested by the kubernetes cluster. The most important routine is the namespace collection and so far the
only routine, which changes the health state of the service. This is because the other two routines (Secret- and ConfigMap-Distribution) depend on 
the namespace list. If the namespace list is empty, the other routines will not process the objects and therefore the service does not work. So if the
namespace collection, or the health operator are not running, the service will be restarted.
