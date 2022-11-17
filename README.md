# kubevirt-addon-manager

Open Cluster Management add-on that deploys KubeVirt onto managed clusters.

## License

This project is licensed under the *Apache License 2.0*. A copy of the license can be found in [LICENSE](LICENSE).

## Installing

If you would like to deploy the kubevirt-addon-manager to a cluster, execute:

```shell
make deploy
```

## Usage

On a hub cluster with the kubevirt-addon-manager running, create a ManagedClusterAddOn in the namespace
of the managed cluster you want KubeVirt deployed onto.

Sample ManagedClusterAddOn (replace `managed_cluster_namespace` with the appropriate managed cluster name):

```yaml
apiVersion: addon.open-cluster-management.io/v1alpha1
kind: ManagedClusterAddOn
metadata:
  name: kubevirt-addon
  namespace: <managed_cluster_namespace>
spec: {}
```

## Development

### Running the controller locally pointing to a remote cluster

If you would like to run the kubevirt-addon-manager outside the cluster, execute:

```shell
make run
```
