package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
	"open-cluster-management.io/addon-framework/pkg/agent"
	"open-cluster-management.io/api/client/addon/clientset/versioned"
)

func Run(ctx context.Context, kubeConfig *rest.Config) error {
	client, err := versioned.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}

	agentAddon, err := addonfactory.NewAgentAddonFactory(addonName, fs, "manifests").
		WithScheme(scheme()).
		WithConfigGVRs(addonfactory.AddOnDeploymentConfigGVR).
		WithGetValuesFuncs(
			defaultValues,
			addonfactory.GetAddOnDeloymentConfigValues(
				addonfactory.NewAddOnDeloymentConfigGetter(client),
				addonfactory.ToAddOnDeloymentConfigValues,
			),
		).
		WithAgentRegistrationOption(registrationOption(
			kubeConfig,
			addonName,
			rand.String(5),
		)).
		WithInstallStrategy(agent.InstallByLabelStrategy(defaultInstallNamespace, v1.LabelSelector{
			MatchLabels: map[string]string{
				managedClusterInstallAddonLabel: managedClusterInstallAddonLabelValue,
			},
		})).
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
