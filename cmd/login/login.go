package login

import (
	"bufio"
	"errors"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var ErrLoginFailed = errors.New("login failed, please check your token")

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [token]",
		Short: "Login to selefra cloud using token",
		Long:  "Login to selefra cloud using token",
		RunE:  RunFunc,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func RunFunc(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		err = MustLogin(args[0])
	}

	err = MustLogin("")

	return err
}

// ShouldLogin should login to selefra cloud
// if login successfully, global token will be set, else return an error
func ShouldLogin(token string) error {
	var err error

	if token == "" {
		token, err = utils.GetCredentialsToken()
		if err != nil {
			ui.PrintErrorLn(err.Error())
			return err
		}
	}

	res, err := httpClient.Login(token)
	if err != nil {
		return ErrLoginFailed
	}
	displayLoginSuccess(res.Data.OrgName, res.Data.TokenName, token)

	global.LOGINTOKEN = token

	return nil
}

// MustLogin unless the user enters wrong token, login is guaranteed
func MustLogin(token string) error {
	var err error

	if err := ShouldLogin(token); err == nil {
		return nil
	}

	token, err = getInputToken()
	if err != nil {
		return errors.New("input token failed")
	}
	if err = ShouldLogin(token); err == nil {
		return nil
	}

	return ErrLoginFailed
}

func getInputToken() (string, error) {
	credentialPath, err := utils.GetCredentialsPath()
	if err != nil {
		return "", err
	}
	ui.PrintCustomizeFNotN(ui.InfoColor, `
Selefra will login for login app.selefra.io  using your browser.
If login is successful, Terraform will store the token in plain text in
the following file for use by subsequent commands:
	%s

	Enter your access token from https://app.selefra.io/settings/access_tokens
	or hit <ENTER> to log in using your browser:`, credentialPath)
	reader := bufio.NewReader(os.Stdin)
	rawToken, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	token := strings.TrimSpace(strings.Replace(rawToken, "\n", "", -1))
	if token == "" {
		ui.PrintErrorLn("No token provided")
		return "", errors.New("no token provided")
	}

	return token, nil
}

func displayLoginSuccess(orgName, tokenName, token string) {
	err := utils.SetCredentials(token)
	if global.LOGINTOKEN == "" {
		global.LOGINTOKEN = token
	}
	global.ORGNAME = orgName
	if err != nil {
		ui.PrintErrorLn(err.Error())
		return
	}
	ui.PrintSuccessF(`
Retrieved token for user: %s. 

Welcome to Selefra Cloud!

Logged in to selefra as %s (https://app.selefra.io/%s)`, tokenName, orgName, orgName)
}
