package login

import (
	"errors"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	"github.com/selefra/selefra/pkg/logger"
	"github.com/spf13/cobra"
)

var ErrLoginFailed = errors.New("login failed, please check your token")

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "login [token]",
		Short:            "Login to selefra cloud using token",
		Long:             "Login to selefra cloud using token",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE:             RunFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func RunFunc(cmd *cobra.Command, args []string) error {

	cli_runtime.Init("./")

	diagnostics := schema.NewDiagnostics()

	host, d := cli_runtime.FindServerHost()
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return err
	}
	logger.InfoF("use server address: %s", host)

	client, d := cloud_sdk.NewCloudClient(host)
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return err
	}
	logger.InfoF("create cloud client success")

	// If you are already logged in, repeat login is not allowed and you must log out first
	getCredentials, _ := client.GetCredentials()
	if getCredentials != nil {
		cli_ui.Errorf("You already logged in as %s, please logout first.\n", getCredentials.UserName)
		return nil
	}

	// Read the token from standard input
	token, d := cli_ui.InputCloudToken(host)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return err
	}

	credentials, d := client.Login(token)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		cli_ui.ShowLoginFailed(token)
		return nil
	}

	cli_ui.ShowLoginSuccess(host, credentials)

	return nil
}
