package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func newCmdProviderInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "install",
		Short:            "Install one or more plugins",
		Long:             "Install one or more plugins",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			err := install(ctx, args)
			return err
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func install(ctx context.Context, args []string) error {
	configYaml, err := config.GetConfig()
	if err != nil {
		ui.Errorln(err.Error())
		return err
	}

	namespace, _, err := utils.Home()
	if err != nil {
		ui.Errorln(err.Error())
		return nil
	}

	provider := registry.NewProviderRegistry(namespace)
	for _, s := range args {
		splitArr := strings.Split(s, "@")
		var name string
		var version string
		if len(splitArr) > 1 {
			name = splitArr[0]
			version = splitArr[1]
		} else {
			name = splitArr[0]
			version = "latest"
		}
		pr := registry.Provider{
			Name:    name,
			Version: version,
			Source:  "",
		}
		p, err := provider.Download(ctx, pr, true)
		continueFlag := false
		for _, provider := range configYaml.Selefra.Providers {
			providerName := *provider.Source
			if strings.ToLower(providerName) == strings.ToLower(p.Name) && strings.ToLower(provider.Version) == strings.ToLower(p.Version) {
				continueFlag = true
				break
			}
		}
		if continueFlag {
			ui.Warningln(fmt.Sprintf("Provider %s@%s already installed", p.Name, p.Version))
			continue
		}
		if err != nil {
			ui.Errorf("Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
			return nil
		} else {
			ui.Successf("Installed %s@%s verified", p.Name, p.Version)
		}
		ui.Infof("Synchronization %s@%s's config...", p.Name, p.Version)
		plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
		if err != nil {
			ui.Errorf("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
		}

		plugProvider := plug.Provider()
		storageOpt := pgstorage.DefaultPgStorageOpts()
		opt, err := json.Marshal(storageOpt)
		initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
			Workspace: utils.ToStringPointer(global.WorkSpace()),
			Storage: &shard.Storage{
				Type:           0,
				StorageOptions: opt,
			},
			IsInstallInit:  pointer.TruePointer(),
			ProviderConfig: pointer.ToStringPointer(""),
		})

		if err != nil {
			ui.Errorln(err.Error())
			return nil
		}

		if initRes != nil && initRes.Diagnostics != nil {
			err := ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
			if err != nil {
				return nil
			}
		}

		res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
		if err != nil {
			ui.Errorf("Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
		}
		ui.Successf("Synchronization %s@%s's config successful", p.Name, p.Version)
		err = tools.SetSelefraProvider(p, configYaml, version)
		if err != nil {
			ui.Errorln(err.Error())
			return nil
		}
		hasProvider := false
		for _, Node := range configYaml.Providers.Content {
			if Node.Kind == yaml.ScalarNode && Node.Value == p.Name {
				hasProvider = true
				break
			}
		}
		if !hasProvider {
			err = tools.SetProviders(res.DefaultConfigTemplate, p, configYaml)
		}
		if err != nil {
			ui.Errorf("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
			return nil
		}
	}

	str, err := yaml.Marshal(configYaml)
	if err != nil {
		ui.Errorln(err.Error())
		return nil
	}
	path, err := config.GetConfigPath()
	if err != nil {
		ui.Errorln(err.Error())
		return nil
	}
	err = os.WriteFile(path, str, 0644)
	return nil
}
