package executors

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/json_util"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"os"
	"path/filepath"
	"sync"
)

// ------------------------------------------------- --------------------------------------------------------------------

// RuleQueryResult 表示一条规则的查询结果
type RuleQueryResult struct {

	// 当前任务的索引编号
	Index int

	// 这个规则是属于哪个模块的
	Module *module.Module

	// 查询出来之后渲染的值是啥
	RuleBlock *module.RuleBlock

	// 查询规则时使用的查询计划是啥，上面会有一些上下文信息之类的
	RulePlan *planner.RulePlan

	// 使用的是哪个provider的哪个版本
	Provider *registry.Provider

	// 使用的配置是哪个
	ProviderConfiguration *module.ProviderBlock

	// 查询的数据库是哪个
	Schema string

	// 查出issue的那一行数据
	Row *schema.Row
}

// ------------------------------------------------- --------------------------------------------------------------------

// ModuleQueryExecutorOptions Option to perform module queries
type ModuleQueryExecutorOptions struct {

	// Query plan to execute
	Plan *planner.ModulePlan

	// The path to install to
	DownloadWorkspace string

	// Receive real-time message feedback
	MessageChannel *message.Channel[*schema.Diagnostics]

	// 执行查询时检测到的rule都往这个channel中放
	RuleQueryResultChannel *message.Channel[*RuleQueryResult]

	// Tracking installation progress
	ProgressTracker getter.ProgressTracker

	// Used to communicate with the provider
	ProviderInformationMap map[string]*shard.GetProviderInformationResponse

	// 每个Provider可能会有多个Fetch任务，只要策略绑定到了这个Provider上，那么Provider所有的Storage都要执行一遍策略
	ProviderExpandMap map[string][]*planner.ProviderContext

	// 查询时使用的并发数
	WorkerNum int
}

// ------------------------------------------------- --------------------------------------------------------------------

const ModuleQueryExecutorName = "module-query-executor"

type ModuleQueryExecutor struct {
	options *ModuleQueryExecutorOptions

	//ruleMetricCounter *RuleMetricCounter
	//ruleMetricChannel chan *RuleMetric
}

var _ Executor = &ModuleQueryExecutor{}

func NewModuleQueryExecutor(options *ModuleQueryExecutorOptions) *ModuleQueryExecutor {
	return &ModuleQueryExecutor{
		options: options,
		//ruleMetricCounter: NewRuleMetricCounter(),
		//ruleMetricChannel: make(chan *RuleMetric, 100),
	}
}

func (x *ModuleQueryExecutor) Name() string {
	return ModuleQueryExecutorName
}

// ------------------------------------------------- --------------------------------------------------------------------

//func (x *ModuleQueryExecutor) StartMetricWorker() {
//	go func() {
//		for metric := range x.ruleMetricChannel {
//			x.ruleMetricCounter.Submit(metric)
//		}
//	}()
//}
//
//func (x *ModuleQueryExecutor) SubmitRuleMetric(rule string, hits int) {
//	x.ruleMetricChannel <- &RuleMetric{Rule: rule, HitCount: hits}
//}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *ModuleQueryExecutor) Execute(ctx context.Context) *schema.Diagnostics {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
		x.options.RuleQueryResultChannel.SenderWaitAndClose()
	}()

	rulePlanSlice := x.makeRulePlanSlice(ctx, x.options.Plan)
	channel := x.toRulePlanChannel(rulePlanSlice)
	x.RunQueryWorker(ctx, channel)

	//close(x.ruleMetricChannel)

	return nil
}

func (x *ModuleQueryExecutor) RunQueryWorker(ctx context.Context, channel chan *planner.RulePlan) {
	wg := sync.WaitGroup{}
	for i := 0; i < x.options.WorkerNum; i++ {
		wg.Add(1)
		NewModuleQueryExecutorWorker(x, channel, &wg).Run(ctx)
	}
	wg.Wait()
}

func (x *ModuleQueryExecutor) toRulePlanChannel(rulePlanSlice []*planner.RulePlan) chan *planner.RulePlan {
	rulePlanChannel := make(chan *planner.RulePlan, len(rulePlanSlice))
	for _, rulePlan := range rulePlanSlice {
		rulePlanChannel <- rulePlan
	}
	close(rulePlanChannel)
	return rulePlanChannel
}

// 把模块及子模块所有的rule执行计划打平，等下要放到一个任务队列中
func (x *ModuleQueryExecutor) makeRulePlanSlice(ctx context.Context, modulePlan *planner.ModulePlan) []*planner.RulePlan {

	rulePlanSlice := make([]*planner.RulePlan, 0)

	// 当前模块的rule执行计划
	rulePlanSlice = append(rulePlanSlice, modulePlan.RulesPlan...)

	// 子模块的执行计划
	for _, subModule := range modulePlan.SubModulesPlan {
		rulePlanSlice = append(rulePlanSlice, x.makeRulePlanSlice(ctx, subModule)...)
	}

	return rulePlanSlice
}

// ------------------------------------------------- --------------------------------------------------------------------

type ModuleQueryExecutorWorker struct {
	ruleChannel chan *planner.RulePlan
	wg          *sync.WaitGroup

	moduleQueryExecutor *ModuleQueryExecutor
}

func NewModuleQueryExecutorWorker(moduleQueryExecutor *ModuleQueryExecutor, rulePlanChannel chan *planner.RulePlan, wg *sync.WaitGroup) *ModuleQueryExecutorWorker {
	return &ModuleQueryExecutorWorker{
		ruleChannel:         rulePlanChannel,
		wg:                  wg,
		moduleQueryExecutor: moduleQueryExecutor,
	}
}

func (x *ModuleQueryExecutorWorker) Run(ctx context.Context) {
	go func() {
		defer func() {
			x.wg.Done()
		}()

		for rulePlan := range x.ruleChannel {
			x.execRulePlan(ctx, rulePlan)
		}

	}()
}

func (x *ModuleQueryExecutorWorker) sendMessage(diagnostics *schema.Diagnostics) {
	if utils.IsNotEmpty(diagnostics) {
		x.moduleQueryExecutor.options.MessageChannel.Send(diagnostics)
	}
}

func (x *ModuleQueryExecutorWorker) execRulePlan(ctx context.Context, rulePlan *planner.RulePlan) {

	x.sendMessage(schema.NewDiagnostics().AddInfo("rule %s begin exec...", rulePlan.String()))

	storages := x.moduleQueryExecutor.options.ProviderExpandMap[rulePlan.BindingProviderName]
	if len(storages) == 0 {
		// TODO 错误报告
		return
	}
	for _, storage := range storages {

		x.execStorageQuery(ctx, rulePlan, storage)
		// TODO 阶段日志
	}
	// TODO 日志

	x.sendMessage(schema.NewDiagnostics().AddInfo("rule %s begin exec done", rulePlan.String()))
}

func (x *ModuleQueryExecutorWorker) execStorageQuery(ctx context.Context, rulePlan *planner.RulePlan, providerContext *planner.ProviderContext) {

	resultSet, diagnostics := providerContext.Storage.Query(ctx, rulePlan.Query)
	if utils.HasError(diagnostics) {
		x.sendMessage(schema.NewDiagnostics().AddErrorMsg("rule %s exec error: %s", rulePlan.String(), diagnostics.ToString()))
		return
	}

	// TODO 打印日志提示
	//x.moduleQueryExecutor.options.MessageChannel <- schema.NewDiagnostics().AddInfo("")
	//cli_ui.Successf("%rootConfig - Rule \"%rootConfig\"\n", rule.Path, rule.Name)
	//cli_ui.Successln("Schema:")
	//cli_ui.Successln(schema + "\n")
	//cli_ui.Successln("Description:")

	for resultSet.Next() {
		rows, d := resultSet.ReadRows(10)
		if rows != nil {
			for _, row := range rows.SplitRowByRow() {
				x.processRuleRow(ctx, rulePlan, providerContext, row)
			}
		}
		if utils.HasError(d) {
			x.sendMessage(d)
		}
	}
}

// 处理rule查询出来的一行
func (x *ModuleQueryExecutorWorker) processRuleRow(ctx context.Context, rulePlan *planner.RulePlan, storage *planner.ProviderContext, row *schema.Row) {
	rowScope := planner.ExtendScope(rulePlan.RuleScope)

	// 把查询出来的行注入到作用域中
	values := row.GetValues()
	for index, columnName := range row.GetColumnNames() {
		rowScope.DeclareVariable(columnName, values[index])
	}

	// 为规则的查询结果渲染出实际的值
	ruleBlockResult, diagnostics := x.renderRule(ctx, rulePlan, rowScope)
	if utils.HasError(diagnostics) {
		x.moduleQueryExecutor.options.MessageChannel.Send(diagnostics)
		return
	}

	result := &RuleQueryResult{
		Module:                rulePlan.Module,
		RulePlan:              rulePlan,
		RuleBlock:             ruleBlockResult,
		Provider:              registry.NewProvider(storage.ProviderName, storage.ProviderVersion),
		ProviderConfiguration: storage.ProviderConfiguration,
		Schema:                storage.Schema,
		Row:                   row,
	}
	x.moduleQueryExecutor.options.RuleQueryResultChannel.Send(result)

	x.sendMessage(schema.NewDiagnostics().AddInfo(json_util.ToJsonString(ruleBlockResult)))

}

func (x *ModuleQueryExecutorWorker) renderRule(ctx context.Context, rulePlan *planner.RulePlan, rowScope *planner.Scope) (*module.RuleBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	ruleBlock := rulePlan.RuleBlock.Copy()

	// 开始渲染相关变量
	// name
	if ruleBlock.Name != "" {
		ruleName, err := rowScope.RenderingTemplate(rulePlan.Name, rulePlan.Name)
		if err != nil {
			// TODO 构造错误上下文
			return nil, diagnostics.AddErrorMsg("render rule name error: %s", err.Error())
		}
		ruleBlock.Name = ruleName
	}

	// labels
	if len(ruleBlock.Labels) > 0 {
		labels := make(map[string]string)
		for key, value := range rulePlan.Labels {
			newValue, err := rowScope.RenderingTemplate(value, value)
			if err != nil {
				// TODO 构造错误上下文
				return nil, diagnostics.AddErrorMsg("render rule labels error: %s", err.Error())
			}
			labels[key] = newValue
		}
		ruleBlock.Labels = labels
	}

	// output
	if ruleBlock.Output != "" {
		output, err := rowScope.RenderingTemplate(rulePlan.Output, rulePlan.Output)
		if err != nil {
			// TODO 构造错误上下文
			return nil, diagnostics.AddErrorMsg("render output labels error: %s", err.Error())
		}
		ruleBlock.Output = output
	}

	// 元数据块的渲染
	d := x.renderRuleMetadata(ctx, rulePlan, ruleBlock, rowScope)
	if diagnostics.AddDiagnostics(d).HasError() {
		return nil, diagnostics
	}

	return ruleBlock, diagnostics
}

// 渲染策略的元数据的块
func (x *ModuleQueryExecutorWorker) renderRuleMetadata(ctx context.Context, rulePlan *planner.RulePlan, ruleBlock *module.RuleBlock, rowScope *planner.Scope) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()
	var err error

	if ruleBlock.MetadataBlock == nil {
		return nil
	}
	metadata := ruleBlock.MetadataBlock

	// description
	if metadata.Description != "" {
		metadata.Description, err = rowScope.RenderingTemplate(metadata.Description, metadata.Description)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("rendering rule description error: %s ", err.Error())
		}
	}

	// title
	if metadata.Title != "" {
		metadata.Title, err = rowScope.RenderingTemplate(metadata.Title, metadata.Title)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("rendering rule title error: %s ", err.Error())
		}
	}

	// 读取修复方案的文本方放上来，如果有必要的话
	if metadata.Remediation != "" {
		markdownFileFullPath := filepath.Join(rulePlan.Module.Workspace, metadata.Remediation)
		file, err := os.ReadFile(markdownFileFullPath)
		if err != nil {
			return diagnostics.AddErrorMsg("read file %s error: %s", markdownFileFullPath, err.Error())
		}
		metadata.Remediation = string(file)
	}

	// tags
	if len(metadata.Tags) != 0 {
		newTags := make([]string, len(metadata.Tags))
		for index, tag := range metadata.Tags {
			newTag, err := rowScope.RenderingTemplate(tag, tag)
			if err != nil {
				// TODO
				return diagnostics.AddErrorMsg("rendering tag error: %s", err.Error())
			}
			newTags[index] = newTag
		}
		metadata.Tags = newTags
	}

	// author
	if metadata.Author != "" {
		author, err := rowScope.RenderingTemplate(metadata.Author, metadata.Author)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("render author error: %s", err.Error())
		}
		metadata.Author = author
	}

	// provider
	if metadata.Provider != "" {
		provider, err := rowScope.RenderingTemplate(metadata.Provider, metadata.Provider)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("render provider error: %s", err.Error())
		}
		metadata.Provider = provider
	}

	// severity
	if metadata.Severity != "" {
		severity, err := rowScope.RenderingTemplate(metadata.Severity, metadata.Severity)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("render severity error: %s", err.Error())
		}
		metadata.Severity = severity
	}

	// id
	if metadata.Id != "" {
		id, err := rowScope.RenderingTemplate(metadata.Id, metadata.Id)
		if err != nil {
			// TODO
			return diagnostics.AddErrorMsg("render id error: %s", err.Error())
		}
		metadata.Id = id
	}

	return diagnostics
}

// ------------------------------------------------- --------------------------------------------------------------------

//type RuleMetricCounter struct {
//	ruleMetricMap map[string]*RuleMetric
//}
//
//func NewRuleMetricCounter() *RuleMetricCounter {
//	return &RuleMetricCounter{
//		ruleMetricMap: make(map[string]*RuleMetric),
//	}
//}
//
//func (x *RuleMetricCounter) Submit(ruleMetric *RuleMetric) {
//	if ruleMetric == nil {
//		return
//	}
//	lastRule, exists := x.ruleMetricMap[ruleMetric.Rule]
//	if !exists {
//		x.ruleMetricMap[ruleMetric.Rule] = ruleMetric
//		return
//	} else {
//		x.ruleMetricMap[ruleMetric.Rule] = ruleMetric.Merge(lastRule)
//	}
//}
//
//// ------------------------------------------------- --------------------------------------------------------------------
//
//type RuleMetric struct {
//	Rule     string
//	HitCount int
//}
//
//func (x *RuleMetric) Merge(other *RuleMetric) *RuleMetric {
//	if x == nil {
//		return other
//	} else if other == nil {
//		return x
//	}
//	if x.Rule != other.Rule {
//		return nil
//	}
//	return &RuleMetric{
//		Rule:     x.Rule,
//		HitCount: x.HitCount + other.HitCount,
//	}
//}

// ------------------------------------------------- --------------------------------------------------------------------

//// 构造表名称到provider名称的映射
//func (x *ModuleQueryExecutor) buildTableToProviderMap() (map[string]string, *schema.Diagnostics) {
//	diagnostics := schema.NewDiagnostics()
//	tableToProviderMap := make(map[string]string, 0)
//	for providerName, providerPlugin := range x.options.ProviderPluginMap {
//		information, err := providerPlugin.Provider().GetProviderInformation(context.Background(), &shard.GetProviderInformationRequest{})
//		if err != nil {
//			return nil, diagnostics
//		}
//		if diagnostics.AddDiagnostics(information.Diagnostics).HasError() {
//			return nil, diagnostics
//		}
//		for tableName := range information.Tables {
//			tableToProviderMap[tableName] = providerName
//		}
//	}
//	return tableToProviderMap, diagnostics
//}

// ------------------------------------------------- --------------------------------------------------------------------
