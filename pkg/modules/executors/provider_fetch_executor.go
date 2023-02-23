package executors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/storage/pgstorage"
	"github.com/selefra/selefra/pkg/utils"
	"io"
	"path/filepath"
	"sync"
	"time"
)

// ------------------------------------------------- --------------------------------------------------------------------

type ProviderFetchExecutorOptions struct {

	// 用来寻找Provider，启动实例
	LocalProviderManager *local_providers_manager.LocalProvidersManager

	// 要执行的拉取计划
	Plans []*planner.ProviderFetchPlan

	// 用于实时接收消息反馈
	MessageChannel *message.Channel[*schema.Diagnostics]

	// 同时拉取的Provider并发数
	WorkerNum int

	// 工作目录
	Workspace string

	// 连接数据库
	DSN string
}

// ------------------------------------------------- --------------------------------------------------------------------

const FetchExecutorName = "fetch-executor"

// ProviderFetchExecutor An actuator for pulling data
type ProviderFetchExecutor struct {
	options *ProviderFetchExecutorOptions

	providerInformationMap map[string]*shard.GetProviderInformationResponse
}

var _ Executor = &ProviderFetchExecutor{}

func NewProviderFetchExecutor(options *ProviderFetchExecutorOptions) *ProviderFetchExecutor {
	return &ProviderFetchExecutor{
		options: options,
	}
}

func (x *ProviderFetchExecutor) GetProviderInformationMap() map[string]*shard.GetProviderInformationResponse {
	return x.providerInformationMap
}

func (x *ProviderFetchExecutor) GetTableToProviderMap() map[string]string {
	tableToProviderMap := make(map[string]string)
	for providerName, providerInformation := range x.providerInformationMap {
		for _, table := range providerInformation.Tables {
			flatTableToProviderMap(providerName, table, tableToProviderMap)
		}
	}
	return tableToProviderMap
}

func flatTableToProviderMap(providerName string, table *schema.Table, m map[string]string) {
	m[table.TableName] = providerName

	for _, subTable := range table.SubTables {
		flatTableToProviderMap(providerName, subTable, m)
	}
}

func (x *ProviderFetchExecutor) Name() string {
	return "provider-fetch-executor"
}

func (x *ProviderFetchExecutor) Execute(ctx context.Context) *schema.Diagnostics {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
	}()

	fetchPlanChannel := make(chan *planner.ProviderFetchPlan, len(x.options.Plans))
	for _, plan := range x.options.Plans {
		fetchPlanChannel <- plan
	}
	close(fetchPlanChannel)

	providerInformationChannel := make(chan *shard.GetProviderInformationResponse, len(x.options.Plans))

	// 开始并发拉取
	wg := sync.WaitGroup{}
	for i := 0; i < x.options.WorkerNum; i++ {
		wg.Add(1)
		NewProviderFetchExecutorWorker(x, fetchPlanChannel, providerInformationChannel, &wg).Run()
	}
	wg.Wait()

	// 整理provider的信息
	close(providerInformationChannel)
	providerInformationMap := make(map[string]*shard.GetProviderInformationResponse)
	for response := range providerInformationChannel {
		providerInformationMap[response.Name] = response
	}
	x.providerInformationMap = providerInformationMap

	return nil
}

// ------------------------------------------------- --------------------------------------------------------------------

// ProviderFetchExecutorWorker 用来执行拉取任务的工作协程
type ProviderFetchExecutorWorker struct {

	// 是执行的哪个执行器中的任务
	executor *ProviderFetchExecutor

	// 任务队列
	planChannel chan *planner.ProviderFetchPlan

	// 退出信号
	wg *sync.WaitGroup

	// 收集启动过的这些Provider的信息
	providerInformation chan *shard.GetProviderInformationResponse
}

func NewProviderFetchExecutorWorker(executor *ProviderFetchExecutor, planChannel chan *planner.ProviderFetchPlan, providerInformation chan *shard.GetProviderInformationResponse, wg *sync.WaitGroup) *ProviderFetchExecutorWorker {
	return &ProviderFetchExecutorWorker{
		executor:            executor,
		planChannel:         planChannel,
		wg:                  wg,
		providerInformation: providerInformation,
	}
}

func (x *ProviderFetchExecutorWorker) Run() {
	go func() {
		defer func() {
			x.wg.Done()
		}()
		for plan := range x.planChannel {
			// 单个Provider的拉取时间限制为24个小时，如果不够的话再调整
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Hour*24)
			x.executePlan(ctx, plan)
			cancelFunc()
		}
	}()
}

// 执行一个provider fetch任务计划
func (x *ProviderFetchExecutorWorker) executePlan(ctx context.Context, plan *planner.ProviderFetchPlan) {

	diagnostics := schema.NewDiagnostics()

	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("begin fetch provider %s", plan.String())))

	// 寻找provider在本地的路径
	localProvider := &local_providers_manager.LocalProvider{
		Provider: plan.Provider,
	}
	installed, d := x.executor.options.LocalProviderManager.IsProviderInstalled(ctx, localProvider)
	if diagnostics.AddDiagnostics(d).HasError() {
		x.sendMessage(x.addProviderNameForMessage(plan, diagnostics))
		return
	}
	if !installed {
		x.sendMessage(x.addProviderNameForMessage(plan, diagnostics.AddErrorMsg("provider %s not installed, can not exec fetch for it", plan.String())))
		return
	}

	// 查找provider在本地安装的位置
	localProviderMeta, d := x.executor.options.LocalProviderManager.Get(ctx, localProvider)
	if diagnostics.AddDiagnostics(d).HasError() {
		x.sendMessage(x.addProviderNameForMessage(plan, diagnostics))
		return
	}

	// 启动provider
	plug, err := plugin.NewManagedPlugin(localProviderMeta.ExecutableFilePath, plan.Name, plan.Version, "", nil)
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("start provider %s at %s failed: %s", plan.String(), localProvider.ExecutableFilePath, err.Error())))
		return
	}
	// 方法执行结束的时候关闭provider退出
	defer plug.Close()

	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("start provider %s success", plan.String())))

	// 数据库连接选项
	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(x.executor.options.DSN)
	dbSchema := pgstorage.GetSchemaKey(plan.Name, plan.Version, plan.ProviderConfigurationBlock)
	pgstorage.WithSearchPath(dbSchema)(storageOpt)
	opt, err := json.Marshal(storageOpt)
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("json marshal postgresql options error: %s", err.Error())))
		return
	}

	// 先获取到锁
	storage, d := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, storageOpt)
	x.sendMessage(x.addProviderNameForMessage(plan, d))
	if utils.HasError(d) {
		return
	}
	lockId := "selefra-fetch-lock"
	ownerId := utils.BuildOwnerId()
	tryTimes := 0
	for {

		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, try get fetch lock...", plan.String(), dbSchema, ownerId)))

		tryTimes++
		err := storage.Lock(ctx, lockId, ownerId)
		if err != nil {
			x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, owner = %s, get fetch lock error: %s, will sleep & retry, tryTimes = %d", plan.String(), dbSchema, ownerId, err.Error(), tryTimes)))
		} else {
			x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, get fetch lock success", plan.String(), dbSchema, ownerId)))
			break
		}
		time.Sleep(time.Second * 10)
	}
	defer func() {
		for tryTimes := 0; tryTimes < 10; tryTimes++ {
			err := storage.UnLock(ctx, lockId, ownerId)
			if err != nil {
				x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, owner = %s, release fetch lock error: %s, will sleep & retry, tryTimes = %d", plan.String(), dbSchema, ownerId, err.Error(), tryTimes)))
			} else {
				x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s, schema %s, owner = %s, release fetch lock success", plan.String(), dbSchema, ownerId)))
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

	workspace, _ := filepath.Abs(x.executor.options.Workspace)
	providerInitResponse, err := pluginProvider.Init(ctx, &shard.ProviderInitRequest{
		Workspace: pointer.ToStringPointer(workspace),
		Storage: &shard.Storage{
			Type:           0,
			StorageOptions: opt,
		},
		IsInstallInit:  pointer.FalsePointer(),
		ProviderConfig: pointer.ToStringPointerOrNilIfEmpty(providerYamlConfiguration),
	})
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("start provider failed: %s", err.Error())))
		return
	}
	if utils.IsNotEmpty(providerInitResponse.Diagnostics) {
		x.sendMessage(x.addProviderNameForMessage(plan, providerInitResponse.Diagnostics))
		if utils.HasError(providerInitResponse.Diagnostics) {
			return
		}
	}
	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s init success", plan.String())))

	// 获取启动的这个provider的相关信息
	information, err := pluginProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, get provider information failed: %s", plan.String(), dbSchema, err.Error())))
		return
	}
	x.providerInformation <- information

	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("get provider %s information success", plan.String())))

	// 删除provider之前的表
	dropRes, err := pluginProvider.DropTableAll(ctx, &shard.ProviderDropTableAllRequest{})
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, drop all table failed: %s", plan.String(), dbSchema, err.Error())))
		return
	}
	x.sendMessage(x.addProviderNameForMessage(plan, dropRes.Diagnostics))
	if utils.HasError(dropRes.Diagnostics) {
		return
	}
	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s drop database schema clean success", plan.String())))

	// 创建所有的表
	createRes, err := pluginProvider.CreateAllTables(ctx, &shard.ProviderCreateAllTablesRequest{})
	if err != nil {
		cli_ui.Errorln(err.Error())
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, create all table failed: %s", plan.String(), dbSchema, err.Error())))
		return
	}
	if createRes.Diagnostics != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, createRes.Diagnostics))
		if utils.HasError(createRes.Diagnostics) {
			return
		}
	}

	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s create tables success", plan.String())))
	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s begin fetch...", plan.String())))

	// 开始拉取数据
	recv, err := pluginProvider.PullTables(ctx, &shard.PullTablesRequest{
		Tables:        plan.GetNeedPullTablesName(),
		MaxGoroutines: plan.GetMaxGoroutines(),
		Timeout:       0,
	})
	if err != nil {
		x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg("provider %s, schema %s, pull table failed: %s", plan.String(), dbSchema, err.Error())))
		return
	}
	//progbar := progress.DefaultProgress()
	//progbar.Add(decl.Name+"@"+decl.Version, -1)
	//success := 0
	//errorsN := 0
	//var total int64
	//for {
	//	res, err := recv.Recv()
	//	if err != nil {
	//		if errors.Is(err, io.EOF) {
	//			progbar.Current(decl.Name+"@"+decl.Version, total, "Done")
	//			progbar.Done(decl.Name + "@" + decl.Version)
	//			break
	//		}
	//		return err
	//	}
	//	progbar.SetTotal(decl.Name+"@"+decl.Version, int64(res.TableCount))
	//	progbar.Current(decl.Name+"@"+decl.Version, int64(len(res.FinishedTables)), res.Table)
	//	total = int64(res.TableCount)
	//	if res.Diagnostics != nil {
	//		if res.Diagnostics.HasError() {
	//			cli_ui.SaveLogToDiagnostic(res.Diagnostics.GetDiagnosticSlice())
	//		}
	//	}
	//	success = len(res.FinishedTables)
	//	errorsN = 0
	//}
	//progbar.ReceiverWait(decl.Name + "@" + decl.Version)
	//if errorsN > 0 {
	//	cli_ui.Errorf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
	//	return nil
	//}
	//cli_ui.Successf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
	//return nil

	success := 0
	errorsN := 0
	var total int64
	for {
		res, err := recv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddErrorMsg(err.Error())))
			return
		}
		//progbar.SetTotal(decl.Name+"@"+decl.Version, int64(res.TableCount))
		//progbar.Current(decl.Name+"@"+decl.Version, int64(len(res.FinishedTables)), res.Table)
		total = int64(res.TableCount)
		if res.Diagnostics != nil {
			//if res.Diagnostics.HasError() {
			//	cli_ui.SaveLogToDiagnostic(res.Diagnostics.GetDiagnosticSlice())
			//}
			x.sendMessage(x.addProviderNameForMessage(plan, res.Diagnostics))
		}
		success = len(res.FinishedTables)
		errorsN = 0
	}
	_ = success
	_ = total
	//progbar.ReceiverWait(decl.Name + "@" + decl.Version)
	if errorsN > 0 {
		//cli_ui.Errorf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
		//return nil
		return
	}
	//cli_ui.Successf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
	//return nil
	x.sendMessage(x.addProviderNameForMessage(plan, schema.NewDiagnostics().AddInfo("provider %s fetch done", plan.String())))

	return
}

func (x *ProviderFetchExecutorWorker) addProviderNameForMessage(plan *planner.ProviderFetchPlan, d *schema.Diagnostics) *schema.Diagnostics {
	if d == nil {
		return nil
	}
	diagnostics := schema.NewDiagnostics()
	for _, item := range d.GetDiagnosticSlice() {
		diagnostics.AddDiagnostic(schema.NewDiagnostic(item.Level(), fmt.Sprintf("provider %s: %s", plan.String(), item.Content())))
	}
	return diagnostics
}

func (x *ProviderFetchExecutorWorker) sendMessage(message *schema.Diagnostics) {
	x.executor.options.MessageChannel.Send(message)
}

// ------------------------------------------------- --------------------------------------------------------------------
