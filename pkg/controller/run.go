package controller

import (
	"context"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager"
)

const (
	addonName = "kubevirt"
)

func Run(ctx context.Context, kubeConfig *rest.Config) error {
	mgr, err := addonmanager.New(kubeConfig)
	if err != nil {
		return err
	}

	agentAddon, err := addonfactory.NewAgentAddonFactory(addonName, fs, "manifests").
		WithScheme(scheme()).
		WithGetValuesFuncs(getDefaultValues).
		WithAgentRegistrationOption(registrationOption(
			kubeConfig,
			addonName,
			rand.String(5),
		)).
		WithInstallStrategy(installStrategy()).
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
