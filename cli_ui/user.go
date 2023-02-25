package cli_ui

import (
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	"github.com/selefra/selefra/pkg/utils"
	"strings"
)

// ------------------------------------------------- --------------------------------------------------------------------

// CloudTokenRequestPath What is the request path to obtain the cloud token
// If there is a change in the address of the cloud side, synchronize it here
const CloudTokenRequestPath = "/Settings/accessTokens"

// InputCloudToken Guide the user to enter a cloud token
func InputCloudToken(serverUrl string) (string, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	if !strings.HasPrefix(strings.ToLower(serverUrl), "http") {
		serverUrl = "https://" + serverUrl
	}

	tipsTemplate := `
selefra will login for login {{.ServerUrl}} using your browser.
if login is successful, selefra will store the token in plain text in
the following file for use by subsequent commands:

	Enter your access token from {{.ServerUrl}}{{.CloudTokenRequestPath}}
	or hit <ENTER> to log in using your browser:`

	// Render display tips
	data := make(map[string]string)
	data["ServerUrl"] = serverUrl
	data["CloudTokenRequestPath"] = CloudTokenRequestPath
	inputCloudTokenTips, err := utils.RenderingTemplate("input-token-tips-template", tipsTemplate, data)
	if err != nil {
		return "", diagnostics.AddErrorMsg("input-token-tips-template render error: %s", err.Error())
	}
	fmt.Println(inputCloudTokenTips)

	// Open a browser window to allow the user to log in
	utils.OpenBrowser(serverUrl)

	// Read the token entered by the user
	var rawToken string
	_, err = fmt.Scanln(&rawToken)
	//reader := bufio.NewReader(os.Stdin)
	//rawToken, err := reader.ReadString('\n')
	if err != nil {
		return "", diagnostics.AddErrorMsg("input cloud token error: %s", err.Error())
	}
	cloudToken := strings.TrimSpace(strings.Replace(rawToken, "\n", "", -1))
	if cloudToken == "" {
		return "", diagnostics.AddErrorMsg("No token provided")
	}

	return cloudToken, diagnostics
}

// ShowLoginSuccess The CLI prompt indicating successful login is displayed
func ShowLoginSuccess(serverUrl string, cloudCredentials *cloud_sdk.CloudCredentials) {
	loginSuccessTemplate := `
Retrieved token for user: {{.UserName}}.

Welcome to Selefra CloudClient!

Logged in to selefra as {{.UserName}} (https://{{.ServerHost}}/{{.OrgName}})
`
	template, err := utils.RenderingTemplate("login-success-tips-template", loginSuccessTemplate, cloudCredentials)
	if err != nil {
		Errorf("render login success message error: %s\n", err.Error())
		return
	}
	Successf(template)
}

// ShowLoginFailed Displays a login failure message
func ShowLoginFailed(cloudToken string) {
	Errorf("You input token %s login failed\n", cloudToken)
}

// ------------------------------------------------- --------------------------------------------------------------------

// ShowRetrievedCloudCredentials Displays the results of the local retrieval of login credentials
func ShowRetrievedCloudCredentials(cloudCredentials *cloud_sdk.CloudCredentials) {
	if cloudCredentials == nil {
		return
	}
	Successf(fmt.Sprintf("Auto login with user %s\n", cloudCredentials.UserName))
}

// ------------------------------------------------- --------------------------------------------------------------------

// ShowLogout Display the logout success prompt
func ShowLogout(cloudCredentials *cloud_sdk.CloudCredentials) {
	if cloudCredentials == nil {
		return
	}
	Successf(fmt.Sprintf("User %s logout success\n", cloudCredentials.UserName))
}

// ------------------------------------------------- --------------------------------------------------------------------
