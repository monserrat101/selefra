package provider

import (
	"fmt"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/version"
	"github.com/spf13/cobra"
)

func newCmdProviderList() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "ListProviders currently installed plugins",
		Long:             "ListProviders currently installed plugins",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {

			downloadWorkspace, err := config.GetDefaultDownloadCacheDirectory()
			if err != nil {
				return err
			}

			return list(downloadWorkspace)
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func list(downloadWorkspace string) error {

	manager, err := local_providers_manager.NewLocalProvidersManager(downloadWorkspace)
	if err != nil {
		return err
	}
	providers, diagnostics := manager.ListProviders()
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return err
	}
	fmt.Printf("  %-13s %-26s %s\n", "Name", "Version", "Source")
	for _, provider := range providers {
		versions := make([]string, 0)
		for versionString := range provider.ProviderVersionMap {
			versions = append(versions, versionString)
		}
		version.Sort(versions)
		for _, versionString := range versions {
			fmt.Printf("  %-13s %-26s %s\n", provider.ProviderName, versionString, provider.ProviderVersionMap[versionString].ExecutableFilePath)
		}
	}
	return nil
}
