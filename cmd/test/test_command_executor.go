package test

import (
	"context"
	"encoding/json"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/module_loader"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/storage/pgstorage"
	"github.com/selefra/selefra/pkg/utils"
	"time"
)

// TestCommandExecutorOptions Used to verify the validity of the module's working directory and configuration
type TestCommandExecutorOptions struct {
	ProjectWorkspace  string
	DownloadWorkspace string
	MessageChannel    *message.Channel[*schema.Diagnostics]
	DSN               string
}

type TestCommandExecutor struct {
	options *TestCommandExecutorOptions
}

func NewTestCommandExecutor(options *TestCommandExecutorOptions) *TestCommandExecutor {
	return &TestCommandExecutor{
		options: options,
	}
}

func (x *TestCommandExecutor) Run(ctx context.Context) {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
	}()

	// 1. load module
	rootModule, b := x.loadModule(ctx)
	if !b {
		return
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("load module from %s success", x.options.ProjectWorkspace))

	// 2. check module
	validatorContext := module.NewValidatorContext()
	d := rootModule.Check(rootModule, validatorContext)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("check module %s syntax ok", x.options.ProjectWorkspace))

	// 3. make providers fetch plan
	localProviderManager, providersFetchPlan, b := x.makeProvidersFetchPlan(ctx, rootModule)
	if !b {
		return
	}

	// 4. check every one provider
	for _, providerInstallPlan := range providersFetchPlan {
		x.processProviderInstallPlan(ctx, localProviderManager, providerInstallPlan)
	}

	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("RequireProvidersBlock verification completed\n\n"))
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("Profile verification completed\n\n"))

}

// load project workspace module in memory
func (x *TestCommandExecutor) loadModule(ctx context.Context) (*module.Module, bool) {
	loader, err := module_loader.NewLocalDirectoryModuleLoader(&module_loader.LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: &module_loader.ModuleLoaderOptions{
			MessageChannel:    x.options.MessageChannel.MakeChildChannel(),
			DownloadDirectory: x.options.DownloadWorkspace,
		},
		ModuleDirectory: x.options.ProjectWorkspace,
	})
	if err != nil {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(err.Error()))
		return nil, false
	}
	return loader.Load(ctx)
}

func (x *TestCommandExecutor) makeProvidersFetchPlan(ctx context.Context, rootModule *module.Module) (*local_providers_manager.LocalProvidersManager, planner.ProvidersFetchPlan, bool) {

	// 1. make install plan
	providersInstallPlan, diagnostics := planner.MakeProviderInstallPlan(ctx, rootModule)
	x.options.MessageChannel.Send(diagnostics)
	if utils.HasError(diagnostics) {
		return nil, nil, false
	}
	if len(providersInstallPlan) == 0 {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("module %s not found providers", x.options.ProjectWorkspace))
		return nil, nil, false
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("module %s find %d providers", x.options.ProjectWorkspace, len(providersInstallPlan)))

	// 2. install providers
	providerInstallExecutor, d := executors.NewProviderInstallExecutor(&executors.ProviderInstallExecutorOptions{
		Plans:             providersInstallPlan,
		MessageChannel:    x.options.MessageChannel.MakeChildChannel(),
		DownloadWorkspace: x.options.DownloadWorkspace,
	})
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return nil, nil, false
	}
	d = providerInstallExecutor.Execute(ctx)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return nil, nil, false
	}

	// 3. make fetch plan
	providersFetchPlan, d := planner.NewProviderFetchPlanner(rootModule, providersInstallPlan.ToMap()).MakePlan(ctx)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return nil, nil, false
	}
	if len(providersFetchPlan) == 0 {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider not fetch plan"))
		return nil, nil, false
	}
	return providerInstallExecutor.GetLocalProviderManager(), providersFetchPlan, true
}

func (x *TestCommandExecutor) processProviderInstallPlan(ctx context.Context, localProviderManager *local_providers_manager.LocalProvidersManager, plan *planner.ProviderFetchPlan) {

	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("begin validate provider %s", plan.String()))

	// Find the local path of the provider
	localProvider := &local_providers_manager.LocalProvider{
		Provider: plan.Provider,
	}
	installed, d := localProviderManager.IsProviderInstalled(ctx, localProvider)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return
	}
	if !installed {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s not installed, can not exec fetch for it", plan.String()))
		return
	}

	// Find the local installation location of the provider
	localProviderMeta, d := localProviderManager.Get(ctx, localProvider)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return
	}

	// start provider
	plug, err := plugin.NewManagedPlugin(localProviderMeta.ExecutableFilePath, plan.Name, plan.Version, "", nil)
	if err != nil {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("start provider %s at %s failed: %s", plan.String(), localProvider.ExecutableFilePath, err.Error()))
		return
	}
	defer plug.Close()

	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("start provider %s success", plan.String()))

	// Database connection option
	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(x.options.DSN)
	dbSchema := pgstorage.GetSchemaKey(plan.Name, plan.Version, plan.ProviderConfigurationBlock)
	pgstorage.WithSearchPath(dbSchema)(storageOpt)
	opt, err := json.Marshal(storageOpt)
	if err != nil {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("json marshal postgresql options error: %s", err.Error()))
		return
	}

	// 先获取到锁
	storage, d := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, storageOpt)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return
	}
	lockId := "selefra-fetch-lock"
	ownerId := utils.BuildOwnerId()
	tryTimes := 0
	for {

		x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, try get fetch lock...", plan.String(), dbSchema, ownerId))

		tryTimes++
		err := storage.Lock(ctx, lockId, ownerId)
		if err != nil {
			x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, owner = %s, get fetch lock error: %s, will sleep & retry, tryTimes = %d", plan.String(), dbSchema, ownerId, err.Error(), tryTimes))
		} else {
			x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, get fetch lock success", plan.String(), dbSchema, ownerId))
			break
		}
		time.Sleep(time.Second * 10)
	}
	defer func() {
		for tryTimes := 0; tryTimes < 10; tryTimes++ {
			err := storage.UnLock(ctx, lockId, ownerId)
			if err != nil {
				x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, owner = %s, release fetch lock error: %s, will sleep & retry, tryTimes = %d", plan.String(), dbSchema, ownerId, err.Error(), tryTimes))
			} else {
				x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, release fetch lock success", plan.String(), dbSchema, ownerId))
				break
			}
		}
	}()

	// 初始化provider
	pluginProvider := plug.Provider()
	var providerYamlConfiguration string
	if plan.ProviderConfigurationBlock == nil {
		providerYamlConfiguration = module.GetDefaultProviderConfigYamlConfiguration(plan.Name, plan.Version)
	} else {
		providerYamlConfiguration = plan.GetProvidersConfigYamlString()
	}

	providerInitResponse, err := pluginProvider.Init(ctx, &shard.ProviderInitRequest{
		Workspace: pointer.ToStringPointer(utils.AbsPath(x.options.ProjectWorkspace)),
		Storage: &shard.Storage{
			Type:           0,
			StorageOptions: opt,
		},
		IsInstallInit:  pointer.FalsePointer(),
		ProviderConfig: pointer.ToStringPointerOrNilIfEmpty(providerYamlConfiguration),
	})
	if err != nil {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("start provider failed: %s", err.Error()))
		return
	}
	if utils.IsNotEmpty(providerInitResponse.Diagnostics) {
		x.options.MessageChannel.Send(providerInitResponse.Diagnostics)
		if utils.HasError(providerInitResponse.Diagnostics) {
			return
		}
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s init success", plan.String()))

	//// 获取启动的这个provider的相关信息
	//information, err := pluginProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
	//if err != nil {
	//	x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, get provider information failed: %s", plan.String(), dbSchema, err.Error()))
	//	return
	//}
	//x.providerInformation <- information
	//x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("get provider %s information success", plan.String()))

	//// 删除provider之前的表
	//dropRes, err := pluginProvider.DropTableAll(ctx, &shard.ProviderDropTableAllRequest{})
	//if err != nil {
	//	x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, drop all table failed: %s", plan.String(), dbSchema, err.Error()))
	//	return
	//}
	//x.options.MessageChannel.Send(dropRes.Diagnostics))
	//if utils.HasError(dropRes.Diagnostics) {
	//	return
	//}
	//x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s drop database schema clean success", plan.String()))
	//
	//// 创建所有的表
	//createRes, err := pluginProvider.CreateAllTables(ctx, &shard.ProviderCreateAllTablesRequest{})
	//if err != nil {
	//	cli_ui.Errorln(err.Error())
	//	x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, create all table failed: %s", plan.String(), dbSchema, err.Error()))
	//	return
	//}
	//if createRes.Diagnostics != nil {
	//	x.options.MessageChannel.Send(createRes.Diagnostics))
	//	if utils.HasError(createRes.Diagnostics) {
	//		return
	//	}
	//}

	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s create tables success", plan.String()))
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("provider %s begin fetch...", plan.String()))

}
