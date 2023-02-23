package logout

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	"github.com/selefra/selefra/pkg/logger"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout to selefra cloud",
		Long:  "Logout to selefra cloud",
		RunE:  RunFunc,
	}

	return cmd
}

func RunFunc(cmd *cobra.Command, args []string) error {

	cli_runtime.Init("./")

	diagnostics := schema.NewDiagnostics()

	// Server address
	host, d := cli_runtime.FindServerHost()
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return err
	}
	logger.InfoF("use server address: %s", host)

	client, d := cloud_sdk.NewCloudClient(host)
	if diagnostics.AddDiagnostics(d).HasError() {
		return cli_ui.PrintDiagnostics(diagnostics)
	}
	logger.InfoF("create cloud client success")

	// If you are not logged in, you are not allowed to log out
	credentials, _ := client.GetCredentials()
	if credentials == nil {
		cli_ui.Errorln("You are not login, please login first.")
		return nil
	}
	logger.InfoF("get credentials success")

	// Destroy the local token
	client.SetToken(credentials.Token)
	if err := cli_ui.PrintDiagnostics(client.Logout()); err != nil {
		return err
	}
	cli_ui.ShowLogout(credentials)
	return nil
}
