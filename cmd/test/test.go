package test

import (
	"context"
	"errors"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
	"sync/atomic"
)

// TODO 2023-2-20 15:32:56 Returns a non-zero value if the test fails
func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "test",
		Short:            "Check whether the configuration is valid",
		Long:             "Check whether the configuration is valid",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {

			//projectWorkspace := "./test_data/test_query_module"
			//downloadWorkspace := "./test_download"

			projectWorkspace := "./"
			downloadWorkspace, _ := config.GetDefaultDownloadCacheDirectory()

			return Test(cmd.Context(), projectWorkspace, downloadWorkspace)
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())

	return cmd
}

func Test(ctx context.Context, projectWorkspace, downloadWorkspace string) error {

	hasError := atomic.Bool{}
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
		if utils.HasError(message) {
			hasError.Store(true)
		}
	})
	d := executors.NewProjectLocalLifeCycleExecutor(&executors.ProjectLocalLifeCycleExecutorOptions{
		ProjectWorkspace:                     projectWorkspace,
		DownloadWorkspace:                    downloadWorkspace,
		MessageChannel:                       messageChannel,
		ProjectLifeCycleStep:                 executors.ProjectLifeCycleStepFetch,
		FetchStep:                            executors.FetchStepGetInformation,
		ProjectCloudLifeCycleExecutorOptions: nil,
		DSN:                                  env.GetDatabaseDsn(),
		FetchWorkerNum: 1,
		QueryWorkerNum: 1,
	}).Execute(context.Background())
	messageChannel.ReceiverWait()

	cli_ui.Successf("RequireProvidersBlock verification completed\n\n")
	cli_ui.Successf("Profile verification completed\n\n")

	if utils.IsNotEmpty(d) {
		_ = cli_ui.PrintDiagnostics(d)
		cli_ui.Errorln("Apply failed")
	} else {
		cli_ui.Successln("Apply Done")
	}

	if hasError.Load() {
		return errors.New("Need help? Known on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
	}
	return nil
}
