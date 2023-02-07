package provider

import (
	"fmt"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
)

func newCmdProviderList() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List currently installed plugins",
		Long:             "List currently installed plugins",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := list()
			return err
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func list() error {
	configYaml, err := config.GetConfig()
	if err != nil {
		ui.Errorln("Error:" + err.Error())
		return nil
	}
	fmt.Printf("  %-13s %-26s %s\n", "Name", "Source", "Version")
	for _, provider := range configYaml.Selefra.Providers {
		fmt.Printf("  %-13s %-26s %s\n", provider.Name, *provider.Source, provider.Version)
	}
	return nil
}
