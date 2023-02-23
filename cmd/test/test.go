package test

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
)

// TODO 2023-2-20 15:32:56 Returns a non-zero value if the test fails
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

func testFunc(cmd *cobra.Command, args []string) error {
	cli_runtime.Init("./")
	dsn, _ := cli_runtime.GetDSN()
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	downloadDirectory, _ := config.GetDefaultDownloadCacheDirectory()
	NewTestCommandExecutor(&TestCommandExecutorOptions{
		ProjectWorkspace:  "./",
		DownloadWorkspace: downloadDirectory,
		MessageChannel:    messageChannel,
		DSN:               dsn,
	}).Run(context.Background())
	messageChannel.ReceiverWait()
	return nil
}

//func TestFunc(ctx context.Context) error {
//	rootConfig, err := config.GetConfig()
//	if err != nil {
//		cli_ui.Errorln("GetWDError:" + err.Error())
//	}
//	return CheckSelefraConfig(ctx, rootConfig)
//}
//
//func testFunc(cmd *cobra.Command, args []string) error {
//	ctx := cmd.Context()
//	return TestFunc(ctx)
//}
//
//// CheckSelefraConfig check if config valid
//func CheckSelefraConfig(ctx context.Context, s *config.RootConfig) error {
//	err := s.TestConfigByNode()
//	if err != nil {
//		return err
//	}
//
//	cli_ui.Successf("Client verification completed\n\n")
//	hasError := false
//	for _, p := range s.Selefra.ProviderDecls {
//		if p.Path == "" {
//			p.Path = utils.GetPathBySource(*p.Source, p.Version)
//		}
//		var providersName = *p.Source
//		plug, err := plugin.NewManagedPlugin(p.Path, providersName, p.Version, "", nil)
//		if err != nil {
//			hasError = true
//			cli_ui.Errorf("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
//			continue
//		}
//		confs, err := tools.ProviderConfigStrs(s, p.Name)
//		if err != nil {
//			hasError = true
//			cli_ui.Errorln(err.Error())
//			continue
//		}
//		for _, conf := range confs {
//			var cp config.ProviderBlock
//			err := yaml.Unmarshal([]byte(conf), &cp)
//			if err != nil {
//				hasError = true
//				cli_ui.Errorln(err.Error())
//				continue
//			}
//
//			storageOpt := pgstorage.DefaultPgStorageOpts()
//			opt, err := json.Marshal(storageOpt)
//
//			provider := plug.Provider()
//			initRes, err := provider.Init(ctx, &shard.ProviderInitRequest{
//				Workspace: pointer.ToStringPointer(global.WorkSpace()),
//				Storage: &shard.Storage{
//					Type:           0,
//					StorageOptions: opt,
//				},
//				IsInstallInit:  pointer.FalsePointer(),
//				ProviderConfig: pointer.ToStringPointer(conf),
//			})
//			if err != nil {
//				hasError = true
//				cli_ui.Errorf("%s@%s verification failed ：%s", providersName, p.Version, err.Error())
//				continue
//			} else {
//				if initRes.Diagnostics != nil {
//					err := cli_ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
//					if err != nil {
//						hasError = true
//					}
//					continue
//				}
//			}
//
//			res, err := provider.SetProviderConfig(ctx, &shard.SetProviderConfigRequest{
//				Storage: &shard.Storage{
//					Type:           0,
//					StorageOptions: opt,
//				},
//				ProviderConfig: pointer.ToStringPointer(conf),
//			})
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				hasError = true
//				continue
//			} else {
//				if res.Diagnostics != nil {
//					err := cli_ui.PrintDiagnostic(res.Diagnostics.GetDiagnosticSlice())
//					if err != nil {
//						hasError = true
//					}
//					continue
//				}
//			}
//			cli_ui.Successf("	%s %s@%s check successfully", cp.Name, providersName, p.Version)
//		}
//	}
//
//	cli_ui.Successf("RequireProvidersBlock verification completed\n\n")
//	cli_ui.Successf("Profile verification completed\n\n")
//	if hasError {
//		return errors.New("Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
//	}
//	return nil
//}
