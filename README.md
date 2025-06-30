# Connections API Design Exploration

This is a test bed for exploring design ideas for OpenDataHub Connections API.

## Setup
I'm testing with kind.
```fish
❯ kind version
kind v0.27.0 go1.23.6 linux/amd64

❯ kind create cluster --name=webhook-testing
# output elided

❯ kind get clusters
enabling experimental podman provider
webhook-testing

❯ kubectl config use-context kind-webhook-testing
Switched to context "kind-webhook-testing".

❯ kubectl version 
Client Version: v1.32.3
Kustomize Version: v5.5.0
Server Version: v1.32.2
```

`make build` will build locally, but `make docker-build load-kind` will create a local container image with
`podman`, save it as `image.tar` in the current directory, and then load it into kind.

Run `make deploy` to load the resources into the cluster.

In another terminal:
```fish
❯ kubectl create deployment my-dep --image=busybox -- date
deployment.apps/my-dep created
```

Back in the original terminal:
```fish
❯ kubectl logs -f -n connections-api-system connections-api-controller-manager-d464b758b-6jqws
2025-06-30T01:10:38Z	INFO	controller-runtime.webhook	Registering webhook	{"path": "/bind-connections-to-workloads"}
2025-06-30T01:10:38Z	INFO	setup	starting manager
2025-06-30T01:10:38Z	INFO	starting server	{"name": "health probe", "addr": "[::]:8081"}
2025-06-30T01:10:38Z	INFO	controller-runtime.webhook	Starting webhook server
2025-06-30T01:10:38Z	INFO	setup	disabling http/2
2025-06-30T01:10:38Z	INFO	controller-runtime.certwatcher	Updated current TLS certificate
2025-06-30T01:10:38Z	INFO	controller-runtime.webhook	Serving webhook server	{"host": "", "port": 9443}
2025-06-30T01:10:38Z	INFO	controller-runtime.certwatcher	Starting certificate poll+watcher	{"interval": "10s"}
2025-06-30T01:11:03Z	INFO	admission	Processing: apps/v1, Kind=Deployment	{"object": {"name":"my-dep","namespace":"default"}, "namespace": "default", "name": "my-dep", "resource": {"group":"apps","version":"v1","resource":"deployments"}, "user": "kubernetes-admin", "requestID": "9fbde638-6736-49b5-9f70-c2e7eca27488"}
2025-06-30T01:11:03Z	INFO	admission	Handling request of type *v1.Deployment from kubernetes-admin	{"object": {"name":"my-dep","namespace":"default"}, "namespace": "default", "name": "my-dep", "resource": {"group":"apps","version":"v1","resource":"deployments"}, "user": "kubernetes-admin", "requestID": "9fbde638-6736-49b5-9f70-c2e7eca27488"}```
