package apply

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
)

func NewApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "apply",
		Short:            "Analyze infrastructure",
		Long:             "Analyze infrastructure",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {

			projectWorkspace := "./test_data/test_query_module"
			downloadWorkspace := "./test_download"

			return Apply(cmd.Context(), projectWorkspace, downloadWorkspace)
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

// ------------------------------------------------- --------------------------------------------------------------------

func Apply(ctx context.Context, projectWorkspace, downloadWorkspace string) error {

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	d := executors.NewProjectLocalLifeCycleExecutor(&executors.ProjectLocalLifeCycleExecutorOptions{
		ProjectWorkspace:                     projectWorkspace,
		DownloadWorkspace:                    downloadWorkspace,
		MessageChannel:                       messageChannel,
		ProjectLifeCycleStep:                 executors.ProjectLifeCycleStepQuery,
		FetchStep:                            executors.FetchStepFetch,
		ProjectCloudLifeCycleExecutorOptions: nil,
		DSN:                                  env.GetDatabaseDsn(),
		FetchWorkerNum:                       1,
		QueryWorkerNum:                       1,
	}).Execute(context.Background())
	messageChannel.ReceiverWait()
	if utils.IsNotEmpty(d) {
		_ = cli_ui.PrintDiagnostics(d)
		cli_ui.Errorln("Apply failed")
	} else {
		cli_ui.Errorln("Apply Done")
	}
	return nil
}

// ------------------------------------------------- --------------------------------------------------------------------
