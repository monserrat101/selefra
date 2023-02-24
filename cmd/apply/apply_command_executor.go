package apply

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/grpc/pb/log"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/module_loader"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/utils"
)

// ------------------------------------------------ ---------------------------------------------------------------------

// ApplyCommandExecutorOptions options for create apply command executor
type ApplyCommandExecutorOptions struct {

	// project path
	ProjectWorkspace string

	// download things put where
	DownloadWorkspace string
}

// ------------------------------------------------ ---------------------------------------------------------------------

// ApplyCommandExecutor the executor for exec selefra apply
type ApplyCommandExecutor struct {
	options *ApplyCommandExecutorOptions

	// project module path
	rootModule *module.Module

	// for sync to cloud
	cloudApplyCommandExecutor *CloudApplyCommandExecutor
}

func NewApplyCommandExecutor(options *ApplyCommandExecutorOptions) *ApplyCommandExecutor {
	return &ApplyCommandExecutor{
		options: options,
	}
}

func (x *ApplyCommandExecutor) Run(ctx context.Context) {

	// load module & check
	if !x.loadModule(ctx) {
		return
	}

	// init cloud
	ok := x.initCloudClient(ctx)
	if !ok {
		_ = x.cloudApplyCommandExecutor.UploadLog(ctx, schema.NewDiagnostics().AddErrorMsg("Selefra Cloud init failed, exit."))
		return
	}
	_ = x.cloudApplyCommandExecutor.UploadLog(ctx, schema.NewDiagnostics().AddInfo("Selefra Cloud init success"))

	// validate module is ok
	validatorContext := module.NewValidatorContext()
	d := x.rootModule.Check(x.rootModule, validatorContext)
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return
	}

	// install provider
	providersInstallPlan, b := x.install(ctx)
	if !b {
		x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_INITIALIZING, log.Status_STATUS_FAILED)
		return
	}
	x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_INITIALIZING, log.Status_STATUS_SUCCESS)

	// fetch data
	fetchExecutor, fetchPlans, b := x.fetch(ctx, providersInstallPlan)
	if !b {
		x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_PULL_INFRASTRUCTURE, log.Status_STATUS_FAILED)
		return
	}
	x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_PULL_INFRASTRUCTURE, log.Status_STATUS_SUCCESS)

	// exec query
	if !x.query(ctx, fetchExecutor, fetchPlans) {
		x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_INFRASTRUCTURE_ANALYSIS, log.Status_STATUS_FAILED)
		return
	}

	x.cloudApplyCommandExecutor.ShutdownAndWait(ctx)
	x.cloudApplyCommandExecutor.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_INFRASTRUCTURE_ANALYSIS, log.Status_STATUS_SUCCESS)

	_ = x.cloudApplyCommandExecutor.UploadLog(ctx, schema.NewDiagnostics().AddInfo("Apply done"))
}

// ------------------------------------------------- --------------------------------------------------------------------

// Load the module to be apply
func (x *ApplyCommandExecutor) loadModule(ctx context.Context) bool {
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	loader, err := module_loader.NewLocalDirectoryModuleLoader(&module_loader.LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: &module_loader.ModuleLoaderOptions{
			DownloadDirectory: x.options.DownloadWorkspace,
			ProgressTracker:   nil,
			MessageChannel:    messageChannel,
		},
		ModuleDirectory: x.options.ProjectWorkspace,
	})
	if err != nil {
		cli_ui.Errorln(fmt.Sprintf("create local directory module loader from %s error: %s", x.options.ProjectWorkspace, err.Error()))
		return false
	}
	rootModule, b := loader.Load(ctx)
	messageChannel.ReceiverWait()
	if !b {
		cli_ui.Errorln(fmt.Sprintf("local directory module loader load  %s failed.", x.options.ProjectWorkspace))
		return false
	}

	x.rootModule = rootModule
	return true
}

// ------------------------------------------------- --------------------------------------------------------------------

// install need providers
func (x *ApplyCommandExecutor) install(ctx context.Context) (planner.ProvidersInstallPlan, bool) {
	// Make an installation plan
	providersInstallPlan, diagnostics := planner.MakeProviderInstallPlan(ctx, x.rootModule)
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, diagnostics); err != nil {
		return nil, false
	}
	if len(providersInstallPlan) == 0 {
		_ = x.cloudApplyCommandExecutor.UploadLog(ctx, schema.NewDiagnostics().AddErrorMsg("no providers"))
		return nil, false
	}

	// Installation-dependent dependency
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = x.cloudApplyCommandExecutor.UploadLog(ctx, message)
		}
	})
	executor, diagnostics := executors.NewProviderInstallExecutor(&executors.ProviderInstallExecutorOptions{
		Plans:             providersInstallPlan,
		MessageChannel:    messageChannel,
		DownloadWorkspace: x.options.DownloadWorkspace,
	})
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, diagnostics); err != nil {
		return nil, false
	}
	d := executor.Execute(context.Background())
	messageChannel.ReceiverWait()
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return nil, false
	}
	return providersInstallPlan, true
}

// ------------------------------------------------- --------------------------------------------------------------------

// Start pulling data
func (x *ApplyCommandExecutor) fetch(ctx context.Context, providersInstallPlan planner.ProvidersInstallPlan) (*executors.ProviderFetchExecutor, planner.ProvidersFetchPlan, bool) {

	// Develop a data pull plan
	providerFetchPlans, d := planner.NewProviderFetchPlanner(&planner.ProviderFetchPlannerOptions{
		Module:                       x.rootModule,
		ProviderVersionVoteWinnerMap: providersInstallPlan.ToMap(),
	}).MakePlan(ctx)
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return nil, nil, false
	}

	// Ready to start pulling
	localProviderManager, err := local_providers_manager.NewLocalProvidersManager(x.options.DownloadWorkspace)
	if err != nil {
		_ = x.cloudApplyCommandExecutor.UploadLog(ctx, schema.NewDiagnostics().AddErrorMsg("create local providers manager failed: %s", err.Error()))
		return nil, nil, false
	}
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = x.cloudApplyCommandExecutor.UploadLog(ctx, message)
		}
	})

	fetchExecutor := executors.NewProviderFetchExecutor(&executors.ProviderFetchExecutorOptions{
		LocalProviderManager: localProviderManager,
		Plans:                providerFetchPlans,
		MessageChannel:       messageChannel,
		WorkerNum:            3,
		Workspace:            x.options.ProjectWorkspace,
		DSN:                  x.cloudApplyCommandExecutor.getDSN(ctx, x.rootModule),
	})
	d = fetchExecutor.Execute(context.Background())
	messageChannel.ReceiverWait()
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return nil, nil, false
	}
	return fetchExecutor, providerFetchPlans, true
}

// ------------------------------------------------- --------------------------------------------------------------------

// Start querying the policy and output the query results to the console and upload them to the cloud
func (x *ApplyCommandExecutor) query(ctx context.Context, fetchExecutor *executors.ProviderFetchExecutor, providerFetchPlans planner.ProvidersFetchPlan) bool {
	plan, d := planner.MakeModuleQueryPlan(ctx, &planner.ModulePlannerOptions{
		Module:             x.rootModule,
		TableToProviderMap: fetchExecutor.GetTableToProviderMap(),
	})
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return false
	}
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		_ = x.cloudApplyCommandExecutor.UploadLog(ctx, d)
	})
	resultQueryResultChannel := message.NewChannel[*executors.RuleQueryResult](func(index int, message *executors.RuleQueryResult) {
		x.cloudApplyCommandExecutor.UploadIssue(ctx, message)
	})
	contextMap, d := providerFetchPlans.BuildProviderContextMap(context.Background(), env.GetDatabaseDsn())
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return false
	}
	queryExecutor := executors.NewModuleQueryExecutor(&executors.ModuleQueryExecutorOptions{
		Plan:                   plan,
		DownloadWorkspace:      x.options.DownloadWorkspace,
		MessageChannel:         messageChannel,
		RuleQueryResultChannel: resultQueryResultChannel,
		ProviderInformationMap: fetchExecutor.GetProviderInformationMap(),
		ProviderExpandMap:      contextMap,
		WorkerNum:              50,
	})
	d = queryExecutor.Execute(ctx)
	messageChannel.ReceiverWait()
	resultQueryResultChannel.ReceiverWait()
	if err := x.cloudApplyCommandExecutor.UploadLog(ctx, d); err != nil {
		return false
	}
	return true
}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *ApplyCommandExecutor) initCloudClient(ctx context.Context) bool {
	host, d := cli_runtime.FindServerHost()
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return false
	}

	cloudApplyCommandExecutor := NewCloudApplyCommandExecutor(host)
	x.cloudApplyCommandExecutor = cloudApplyCommandExecutor
	return cloudApplyCommandExecutor.initCloudClient(ctx, x.rootModule)
}

// ------------------------------------------------- --------------------------------------------------------------------
