package main

import (
	goflag "flag"
	"fmt"
	"github.com/0xFelix/kubevirt-addon-manager/pkg/controller"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilflag "k8s.io/component-base/cli/flag"
	logs "k8s.io/component-base/logs/api/v1"
	"open-cluster-management.io/addon-framework/examples/cmdfactory"
	"open-cluster-management.io/addon-framework/pkg/version"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.AddFlags(logs.NewLoggingConfiguration(), pflag.CommandLine)

	command := newCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubevirt-addon-manager",
		Short: "kubevirt-addon-manager",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			os.Exit(1)
		},
	}

	if v := version.Get().String(); len(v) == 0 {
		cmd.Version = "<unknown>"
	} else {
		cmd.Version = v
	}

	cmd.AddCommand(newControllerCommand())

	return cmd
}

func newControllerCommand() *cobra.Command {
	cmd := cmdfactory.
		NewControllerCommandConfig("kubevirt-addon-manager", version.Get(), controller.Run).
		NewCommand()
	cmd.Use = "controller"
	cmd.Short = "Start the kubevirt-addon-manager controller"

	return cmd
}
