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


	//	cli_ui.Successf("RequireProvidersBlock verification completed\n\n")
	//	cli_ui.Successf("Profile verification completed\n\n")
	//	if hasError {
	//		return errors.New("Need help? Know on Slack or open a Github Issue: https://github.com/selefra/selefra#community")
	//	}

	return nil
}
