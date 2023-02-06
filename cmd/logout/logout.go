package logout

import (
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
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
	token, err := utils.GetCredentialsToken()
	if err != nil {
		return err
	}

	return shouldLogout(token)
}

func shouldLogout(token string) error {
	err := httpClient.Logout(token)
	if err != nil {
		ui.Errorln("Logout error:" + err.Error())
		return nil
	}

	err = utils.SetCredentials("")
	if err != nil {
		ui.Errorln(err.Error())
	}

	return nil
}
