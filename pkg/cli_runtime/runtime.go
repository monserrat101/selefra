package cli_runtime

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/cli_env"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/module_loader"
	"github.com/selefra/selefra/pkg/utils"
)

// Runtime 命令行的运行时
var Runtime *CLIRuntime

type CLIRuntime struct {

	// 工作目录是哪个
	Workspace string

	// 下载到哪个目录中
	DownloadWorkspace string

	// 操作时可能会出现的错误
	Diagnostics *schema.Diagnostics

	// 工作目录下的根模块
	RootModule *module.Module

	CloudClient *cloud_sdk.CloudClient
}

func Init(workspace string) {
	Runtime = NewCLIRuntime(workspace)
	Runtime.LoadWorkspaceModule()
}

func NewCLIRuntime(workspace string) *CLIRuntime {
	x := &CLIRuntime{
		Workspace: workspace,
	}
	return x
}

func (x *CLIRuntime) InitCloudClient() {
	host, diagnostics := FindServerHost()
	x.Diagnostics.AddDiagnostics(diagnostics)
	if utils.HasError(diagnostics) {
		return
	}
	client, d := cloud_sdk.NewCloudClient(host)
	x.Diagnostics.AddDiagnostics(d)
	if utils.HasError(d) {
		return
	}
	x.CloudClient = client

	// 如果本地有凭证的话则自动登录
	credentials, _ := client.GetCredentials()
	if credentials != nil {
		login, d := client.Login(credentials.Token)
		if utils.HasError(d) {
			cli_ui.ShowLoginFailed(credentials.Token)
			return
		}
		cli_ui.ShowLoginSuccess(host, login)
	}

}

func (x *CLIRuntime) LoadWorkspaceModule() *CLIRuntime {

	if utils.HasError(x.Diagnostics) {
		return x
	}

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		// TODO log
	})
	loader, err := module_loader.NewLocalDirectoryModuleLoader(&module_loader.LocalDirectoryModuleLoaderOptions{
		ModuleDirectory: x.Workspace,
		ModuleLoaderOptions: &module_loader.ModuleLoaderOptions{
			MessageChannel: messageChannel,
		},
	})
	if err != nil {
		messageChannel.SenderWaitAndClose()
		x.Diagnostics.AddErrorMsg("create module load from directory %s error: %s", x.Workspace, err.Error())
		return x
	}
	workspaceModule, _ := loader.Load(context.Background())
	messageChannel.ReceiverWait()
	if workspaceModule != nil {
		x.RootModule = workspaceModule
	}

	return x
}

// ------------------------------------------------- --------------------------------------------------------------------

const DefaultServerURL = "app.selefra.io"

func FindServerHost() (string, *schema.Diagnostics) {

	// 尝试从配置文件中获取
	if Runtime.RootModule != nil &&
		Runtime.RootModule.SelefraBlock != nil &&
		Runtime.RootModule.SelefraBlock.CloudBlock != nil &&
		Runtime.RootModule.SelefraBlock.CloudBlock.HostName != "" {
		return Runtime.RootModule.SelefraBlock.CloudBlock.HostName, nil
	}

	// 尝试从环境变量中获取
	if cli_env.GetServerHost() != "" {
		return cli_env.GetServerHost(), nil
	}

	// 都获取不到，使用默认的
	return DefaultServerURL, nil
}

// ------------------------------------------------- --------------------------------------------------------------------

func GetDSN() (string, *schema.Diagnostics) {

	// 如果有在当前模块中配置的话优先使用当前模块的配置
	if Runtime != nil && Runtime.RootModule != nil && Runtime.RootModule.SelefraBlock != nil && Runtime.RootModule.SelefraBlock.ConnectionBlock != nil {
		return Runtime.RootModule.SelefraBlock.ConnectionBlock.BuildDSN(), nil
	}

	// 否则看是否登录
	if Runtime.CloudClient != nil && Runtime.CloudClient.IsLoggedIn() {
		return Runtime.CloudClient.FetchOrgDSN()
	}

	// 环境变量
	if env.GetDatabaseDsn() != "" {
		return env.GetDatabaseDsn(), nil
	}

	// TODO 内置的PG数据库
	return "", nil
}

// ------------------------------------------------- --------------------------------------------------------------------
