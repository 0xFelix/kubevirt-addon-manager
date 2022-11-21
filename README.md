# kubevirt-addon-manager

Open Cluster Management add-on that
deploys [KubeVirt](https://github.com/kubevirt/hyperconverged-cluster-operator)
HCO onto managed OpenShift clusters using OLM.

## License

This project is licensed under the *Apache License 2.0*. A copy of the license
can be found in [LICENSE](LICENSE).

## Installing

If you would like to deploy the kubevirt-addon-manager to a cluster, execute:

```shell
make deploy
```

## Usage

### With Label

On a hub cluster with the kubevirt-addon-manager running, add the label
`addons.open-cluster-management.io/kubevirt` with value `true` to the
`ManagedCluster` object of the managed OpenShift cluster you want to deploy
KubeVirt HCO onto.

Sample `ManagedCluster` (replace `managed_cluster_name` with the appropriate
managed cluster name):

```yaml
apiVersion: cluster.open-cluster-management.io/v1
kind: ManagedCluster
metadata:
  name: <managed_cluster_name>
  labels:
    addons.open-cluster-management.io/kubevirt: "true"
spec: ...
```

### Manually creating the ManagedClusterAddOn

On a hub cluster with the kubevirt-addon-manager running, create a
`ManagedClusterAddOn` object in the namespace
of the managed OpenShift cluster you want to deploy KubeVirt HCO onto.

Sample `ManagedClusterAddOn` (replace `managed_cluster_namespace` with the
appropriate managed cluster name):

```yaml
apiVersion: addon.open-cluster-management.io/v1alpha1
kind: ManagedClusterAddOn
metadata:
  name: kubevirt
  namespace: <managed_cluster_namespace>
spec: { }
```

## Development

### Running the controller locally pointing to a remote cluster

If you would like to run the kubevirt-addon-manager controller outside a
cluster, execute:

```shell
make run
```

This will use the kubeconfig found in environment variable `KUBECONFIG` or
default to `~/.kube/config`.
