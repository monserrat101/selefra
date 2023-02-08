package test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "test",
		Short:            "Check whether the configuration is valid",
		Long:             "Check whether the configuration is valid",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE:             testFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())

	return cmd
}

func TestFunc(ctx context.Context) error {
	rootConfig, err := config.GetConfig()
	if err != nil {
		ui.Errorln("GetWDError:" + err.Error())
	}
	return CheckSelefraConfig(ctx, rootConfig)
}

func testFunc(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return TestFunc(ctx)
}

// CheckSelefraConfig check if config valid
func CheckSelefraConfig(ctx context.Context, s *config.RootConfig) error {
	err := s.TestConfigByNode()
	if err != nil {
		return err
	}

	ui.Successf("Client verification completed\n\n")
	hasError := false
	for _, p := range s.Selefra.ProviderDecls {
		if p.Path == "" {
			p.Path = utils.GetPathBySource(*p.Source, p.Version)
		}
		var providersName = *p.Source
		plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
		if err != nil {
			hasError = true
			ui.Errorf("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
			continue
		}
		confs, err := tools.ProviderConfigStrs(s, p.Name)
		if err != nil {
			hasError = true
			ui.Errorln(err.Error())
			continue
		}
		for _, conf := range confs {
			var cp config.Provider
			err := yaml.Unmarshal([]byte(conf), &cp)
			if err != nil {
				hasError = true
				ui.Errorln(err.Error())
				continue
			}

			storageOpt := pgstorage.DefaultPgStorageOpts()
			opt, err := json.Marshal(storageOpt)

			provider := plug.Provider()
			initRes, err := provider.Init(ctx, &shard.ProviderInitRequest{
				Workspace: pointer.ToStringPointer(global.WorkSpace()),
				Storage: &shard.Storage{
					Type:           0,
					StorageOptions: opt,
				},
				IsInstallInit:  pointer.FalsePointer(),
				ProviderConfig: pointer.ToStringPointer(conf),
			})
			if err != nil {
				hasError = true
				ui.Errorf("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
				continue
			} else {
				if initRes.Diagnostics != nil {
					err := ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
					if err != nil {
						hasError = true
					}
					continue
				}
			}

			res, err := provider.SetProviderConfig(ctx, &shard.SetProviderConfigRequest{
				Storage: &shard.Storage{
					Type:           0,
					StorageOptions: opt,
				},
				ProviderConfig: pointer.ToStringPointer(conf),
			})
			if err != nil {
				ui.Errorln(err.Error())
				hasError = true
				continue
			} else {
				if res.Diagnostics != nil {
					err := ui.PrintDiagnostic(res.Diagnostics.GetDiagnosticSlice())
					if err != nil {
						hasError = true
					}
					continue
				}
			}
			ui.Successf("	%s %s@%s check successfully", cp.Name, providersName, p.Version)
		}
	}

	ui.Successf("ProviderDecls verification completed\n\n")
	ui.Successf("Profile verification completed\n\n")
	if hasError {
		return errors.New("Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
	}
	return nil
}
