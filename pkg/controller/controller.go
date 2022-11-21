package controller

import (
	"embed"
	"fmt"

	hcov1beta1 "github.com/kubevirt/hyperconverged-cluster-operator/api/v1beta1"
	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	k8sv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"open-cluster-management.io/addon-framework/examples/rbac"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/addon-framework/pkg/utils"
	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workv1 "open-cluster-management.io/api/work/v1"
)

const (
	managedClusterInstallAddonLabel      = "addons.open-cluster-management.io/kubevirt"
	managedClusterInstallAddonLabelValue = "true"

	defaultInstallNamespace   = "kubevirt-hyperconverged"
	defaultHyperConvergedName = "kubevirt-hyperconverged"
)

//go:embed manifests
var fs embed.FS

func scheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = operatorsv1.AddToScheme(s)
	_ = operatorsv1alpha1.AddToScheme(s)
	_ = hcov1beta1.AddToScheme(s)
	return s
}

func registrationOption(kubeConfig *rest.Config, addonName, agentName string) *agent.RegistrationOption {
	return &agent.RegistrationOption{
		CSRConfigurations: agent.KubeClientSignerConfigurations(addonName, agentName),
		CSRApproveCheck:   utils.DefaultCSRApprover(agentName),
		PermissionConfig:  rbac.AddonRBAC(kubeConfig),
	}
}

func getDefaultValues(cluster *clusterv1.ManagedCluster, addon *addonv1alpha1.ManagedClusterAddOn) (addonfactory.Values, error) {
	manifestConfig := struct {
		ClusterName        string
		HyperConvergedName string
		InstallNamespace   string
		KubeConfigSecret   string
	}{
		ClusterName:        cluster.Name,
		HyperConvergedName: defaultHyperConvergedName,
		InstallNamespace:   defaultInstallNamespace,
		KubeConfigSecret:   fmt.Sprintf("%s-hub-kubeconfig", addon.Name),
	}

	return addonfactory.StructToValues(manifestConfig), nil
}

func installStrategy() *agent.InstallStrategy {
	return agent.InstallByLabelStrategy(defaultInstallNamespace, k8sv1.LabelSelector{
		MatchLabels: map[string]string{
			managedClusterInstallAddonLabel: managedClusterInstallAddonLabelValue,
		},
	})
}

func agentHealthProber() *agent.HealthProber {
	return &agent.HealthProber{
		Type: agent.HealthProberTypeWork,
		WorkProber: &agent.WorkHealthProber{
			ProbeFields: []agent.ProbeField{
				{
					ResourceIdentifier: workv1.ResourceIdentifier{
						Group:     "hco.kubevirt.io",
						Resource:  "hyperconvergeds",
						Name:      defaultHyperConvergedName,
						Namespace: defaultInstallNamespace,
					},
					ProbeRules: []workv1.FeedbackRule{
						{
							Type: workv1.JSONPathsType,
							JsonPaths: []workv1.JsonPath{
								{
									Name: "isAvailable",
									Path: ".status.conditions[?(@.type==\"Available\")].status",
								},
							},
						},
					},
				},
			},
			HealthCheck: func(identifier workv1.ResourceIdentifier, result workv1.StatusFeedbackResult) error {
				if len(result.Values) == 0 {
					return fmt.Errorf("no values are probed for hyperconverged %s/%s", identifier.Namespace, identifier.Name)
				}
				for _, value := range result.Values {
					if value.Name != "isAvailable" {
						continue
					}
					if *value.Value.String == "True" {
						return nil
					}
					return fmt.Errorf("isAvailable is %s for hyperconverged %s/%s", *value.Value.String, identifier.Namespace, identifier.Name)
				}
				return fmt.Errorf("isAvailable is not probed")
			},
		},
	}
}
