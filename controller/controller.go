package controller

import (
	"context"
	"embed"
	"fmt"

	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/examples/rbac"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	"open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/addon-framework/pkg/utils"
	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	addonclient "open-cluster-management.io/api/client/addon/clientset/versioned"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	workv1 "open-cluster-management.io/api/work/v1"
)

const (
	addonName        = "kubevirt-addon"
	installNamespace = "kubevirt-hyperconverged"
)

//go:embed manifests
var fs embed.FS

func newRegistrationOption(kubeConfig *rest.Config, addonName, agentName string) *agent.RegistrationOption {
	return &agent.RegistrationOption{
		CSRConfigurations: agent.KubeClientSignerConfigurations(addonName, agentName),
		CSRApproveCheck:   utils.DefaultCSRApprover(agentName),
		PermissionConfig:  rbac.AddonRBAC(kubeConfig),
	}
}

func getDefaultValues(cluster *clusterv1.ManagedCluster, addon *addonv1alpha1.ManagedClusterAddOn) (addonfactory.Values, error) {
	manifestConfig := struct {
		ClusterName      string
		InstallNamespace string
		KubeConfigSecret string
	}{
		ClusterName:      cluster.Name,
		InstallNamespace: installNamespace,
		KubeConfigSecret: fmt.Sprintf("%s-hub-kubeconfig", addon.Name),
	}

	return addonfactory.StructToValues(manifestConfig), nil
}

func agentHealthProber() *agent.HealthProber {
	return &agent.HealthProber{
		Type: agent.HealthProberTypeWork,
		WorkProber: &agent.WorkHealthProber{
			ProbeFields: []agent.ProbeField{
				{
					ResourceIdentifier: workv1.ResourceIdentifier{
						Group:     "hco.kubevirt.io",
						Resource:  "hyperconverged",
						Name:      "kubevirt-hyperconverged",
						Namespace: installNamespace,
					},
					ProbeRules: []workv1.FeedbackRule{
						{
							Type: workv1.JSONPathsType,
							JsonPaths: []workv1.JsonPath{
								{
									Name: "Available",
									Path: "{.status.conditions[?(@.type == 'Available')].status}",
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
					if value.Name != "Available" {
						continue
					}
					if *value.Value.String == "True" {
						return nil
					}
					return fmt.Errorf("Available is %s for hyperconverged %s/%s", *value.Value.String, identifier.Namespace, identifier.Name)
				}
				return fmt.Errorf("Available is not probed")
			},
		},
	}
}

func Run(ctx context.Context, kubeConfig *rest.Config) error {
	client, err := addonclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}

	registrationOption := newRegistrationOption(
		kubeConfig,
		addonName,
		rand.String(5),
	)

	agentAddon, err := addonfactory.NewAgentAddonFactory(addonName, fs, "manifests").
		WithConfigGVRs(addonfactory.AddOnDeploymentConfigGVR).
		WithGetValuesFuncs(
			getDefaultValues,
			addonfactory.GetAddOnDeloymentConfigValues(
				addonfactory.NewAddOnDeloymentConfigGetter(client),
				addonfactory.ToAddOnDeloymentConfigValues,
			),
		).
		WithAgentRegistrationOption(registrationOption).
		WithInstallStrategy(agent.InstallAllStrategy(installNamespace)).
		WithAgentHealthProber(agentHealthProber()).
		BuildTemplateAgentAddon()
	if err != nil {
		klog.Errorf("failed to build agent %v", err)
		return err
	}

	err = mgr.AddAgent(agentAddon)
	if err != nil {
		klog.Fatal(err)
	}

	err = mgr.Start(ctx)
	if err != nil {
		klog.Fatal(err)
	}
	<-ctx.Done()

	return nil
}
