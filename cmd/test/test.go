package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
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
	err := config.IsSelefra()
	if err != nil {
		ui.Errorln(err.Error())
		return err
	}

	if err != nil {
		ui.Errorln("GetWDError:" + err.Error())
	}
	s := config.RootConfig{}
	return CheckSelefraConfig(ctx, s)
}

func testFunc(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return TestFunc(ctx)
}

func checkConfig(ctx context.Context, c config.RootConfig) error {
	var err error
	if c.Selefra.CliVersion == "" {
		err = errors.New("cliVersion is empty")
		return err
	}
	if c.Selefra.Name == "" {
		err = errors.New("name is empty")
		return err
	}
	uid, _ := uuid.NewUUID()
	for i := range c.Selefra.Providers {
		confs, err := tools.GetProviders(&c, c.Selefra.Providers[i].Name)
		if err != nil {
			ui.Errorln(err.Error())
			return nil
		}
		for _, conf := range confs {
			var cp config.ProviderConfig
			err := yaml.Unmarshal([]byte(conf), &cp)
			if err != nil {
				ui.Errorln(err.Error())
				continue
			}
			_, e := client.CreateClientFromConfig(ctx, &c.Selefra, uid, c.Selefra.Providers[i], cp)
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func CheckSelefraConfig(ctx context.Context, s config.RootConfig) error {
	err := s.TestConfigByNode()
	if err != nil {
		return err
	}
	err = checkConfig(ctx, s)
	if err != nil {
		return errors.New(fmt.Sprintf("selefra configuration exception:%s", err.Error()))
	}
	ui.Successf("Client verification completed")
	hasError := false
	for _, p := range s.Selefra.Providers {
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
		confs, err := tools.GetProviders(&s, p.Name)
		if err != nil {
			hasError = true
			ui.Errorln(err.Error())
			continue
		}
		for _, conf := range confs {
			var cp config.ProviderConfig
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
				Workspace: utils.ToStringPointer(global.WorkSpace()),
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

	ui.Successf("\nProviders verification completed\n")
	ui.Successf("Profile verification completed\n")
	if hasError {
		return errors.New("Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
	}
	return nil
}
