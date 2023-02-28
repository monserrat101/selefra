package executors

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/cli_env"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	selefraGrpc "github.com/selefra/selefra/pkg/grpc"
	"github.com/selefra/selefra/pkg/grpc/pb/issue"
	"github.com/selefra/selefra/pkg/grpc/pb/log"
	"github.com/selefra/selefra/pkg/logger"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/storage/pgstorage"
	"github.com/selefra/selefra/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync/atomic"
)

// ------------------------------------------------- --------------------------------------------------------------------

type ProjectCloudLifeCycleExecutorOptions struct {

	// The address of the Cloud cluster to which you are connecting
	CloudServerHost string

	// The current project is loaded as a module
	Module *module.Module

	// Send messages to the outside world
	MessageChannel *message.Channel[*schema.Diagnostics]

	// Whether to enable console prompts
	EnableConsoleTips bool

	IsNeedLogin bool
}

// ------------------------------------------------- --------------------------------------------------------------------

type ProjectCloudLifeCycleExecutor struct {
	options *ProjectCloudLifeCycleExecutorOptions

	// the client for connect to cloud
	cloudClient *cloud_sdk.CloudClient

	// for upload log
	logClient         log.LogClient
	logStreamUploader *selefraGrpc.StreamUploader[log.Log_UploadLogStreamClient, int, *log.UploadLogStream_Request, *log.UploadLogStream_Response]

	// for upload issue
	issueStreamUploader *selefraGrpc.StreamUploader[issue.Issue_UploadIssueStreamClient, int, *issue.UploadIssueStream_Request, *issue.UploadIssueStream_Response]

	// index generator
	logIdGenerator   atomic.Int64
	issueIdGenerator atomic.Int64

	// task current stage
	stage log.StageType
}

func NewProjectCloudLifeCycleExecutor(options *ProjectCloudLifeCycleExecutorOptions) *ProjectCloudLifeCycleExecutor {
	return &ProjectCloudLifeCycleExecutor{
		options: options,
	}
}

func (x *ProjectCloudLifeCycleExecutor) getServerHost() string {
	if x.options.CloudServerHost != "" {
		return x.options.CloudServerHost
	}
	return cli_env.GetServerHost()
}

func (x *ProjectCloudLifeCycleExecutor) InitCloudClient(ctx context.Context) bool {

	// 1. create cloud client
	cloudClient, d := cloud_sdk.NewCloudClient(x.getServerHost())
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return false
	}
	x.cloudClient = cloudClient

	// 2. find local cloud token & use it to login to the cloud
	if !x.options.IsNeedLogin {
		return true
	}
	credentials, _ := cloudClient.GetCredentials()
	if credentials != nil && !x.loginByCredentials(ctx, credentials) {
		return false
	}

	return true
}

// Login against credentials
func (x *ProjectCloudLifeCycleExecutor) loginByCredentials(ctx context.Context, credentials *cloud_sdk.CloudCredentials) bool {

	if x.cloudClient == nil {
		return false
	}

	// try login
	login, diagnostics := x.cloudClient.Login(credentials.Token)
	x.options.MessageChannel.Send(diagnostics)

	// login failed
	if utils.HasError(diagnostics) {
		if x.options.EnableConsoleTips {
			// login success
			cli_ui.ShowLoginFailed(x.options.CloudServerHost)
		}
		return false
	}

	// login success
	if x.options.EnableConsoleTips {
		cli_ui.ShowLoginSuccess(x.options.CloudServerHost, login)
	}

	// check relative project
	if x.options.Module.SelefraBlock == nil ||
		x.options.Module.SelefraBlock.CloudBlock == nil ||
		x.options.Module.SelefraBlock.CloudBlock.Project == "" {
		errorMsg := fmt.Sprintf("Failed to connect to the cloud, you must specify the project name %s in module", x.options.Module.BuildFullName())
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorMsg))
		return false
	}

	// so, we can get project name now
	projectName := x.options.Module.SelefraBlock.CloudBlock.Project
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo(fmt.Sprintf("try relative module %s to project %s ", x.options.Module.BuildFullName(), projectName)))

	// try to relative project
	project, d := x.cloudClient.CreateProject(projectName)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		return false
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("Successfully connected to cloud, associated module %s to project %s", x.options.Module.BuildFullName(), project.Name))

	// create task
	task, d := x.cloudClient.CreateTask(project.Name)
	x.options.MessageChannel.Send(d)
	if utils.HasError(d) {
		msg := fmt.Sprintf("Failed to create a task for the project %s", project.Name)
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(msg))
		return false
	}

	msg := fmt.Sprintf("Succeeded in creating a task %s for project %s", task.TaskId, project.Name)
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo(msg))

	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("Begin init log & issue uploader..."))
	if !x.initLogUploader(x.cloudClient) {
		return false
	}
	if !x.initIssueUploader(x.cloudClient) {
		return false
	}
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("Init log uploader & issue done"))

	// change task status to begin
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("change task status to INITIALIZING"))

	_ = x.UploadLog(ctx, schema.NewDiagnostics().AddInfo("begin run task %s INITIALIZING stage", task.TaskId))
	return true
}

// init issue uploader for send issue to cloud
func (x *ProjectCloudLifeCycleExecutor) initIssueUploader(client *cloud_sdk.CloudClient) bool {
	issueStreamUploaderMessageChannel := x.options.MessageChannel.MakeChildChannel()
	issueStreamUploader, diagnostics := client.NewIssueStreamUploader(issueStreamUploaderMessageChannel)
	x.options.MessageChannel.Send(diagnostics)
	if utils.HasError(diagnostics) {
		issueStreamUploaderMessageChannel.SenderWaitAndClose()
		return false
	}
	issueStreamUploader.RunUploaderWorker()
	x.issueStreamUploader = issueStreamUploader
	return true
}

// init log uploader for send log to loud
func (x *ProjectCloudLifeCycleExecutor) initLogUploader(client *cloud_sdk.CloudClient) bool {
	logStreamUploaderMessageChannel := x.options.MessageChannel.MakeChildChannel()
	logClient, logStreamUploader, diagnostics := client.NewLogStreamUploader(logStreamUploaderMessageChannel)
	x.options.MessageChannel.Send(diagnostics)
	if utils.HasError(diagnostics) {
		logStreamUploaderMessageChannel.SenderWaitAndClose()
		return false
	}
	logStreamUploader.RunUploaderWorker()
	x.logClient = logClient
	x.logStreamUploader = logStreamUploader
	return true
}

// ------------------------------------------------- --------------------------------------------------------------------

// UploadIssue add issue to send cloud queue
func (x *ProjectCloudLifeCycleExecutor) UploadIssue(ctx context.Context, r *RuleQueryResult) {
	var consoleOutput strings.Builder
	consoleOutput.WriteString(fmt.Sprintf("Rule name %s, ", r.RuleBlock.Name))
	if r.RuleBlock.MetadataBlock != nil && r.RuleBlock.MetadataBlock.Id != "" {
		consoleOutput.WriteString(fmt.Sprintf("id %s, ", r.RuleBlock.MetadataBlock.Id))
	}
	consoleOutput.WriteString(fmt.Sprintf("output %s", r.RuleBlock.Output))
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo(consoleOutput.String()))

	// send to cloud
	if x.issueStreamUploader == nil {
		logger.ErrorF("issueStreamUploader is nil, ignore issue upload")
		return
	}
	request := x.convertRuleQueryResultToIssueUploadRequest(r)
	x.issueStreamUploader.Submit(ctx, int(request.Index), request)
}

// Convert the query results of the rules into a format uploaded to the Cloud
func (x *ProjectCloudLifeCycleExecutor) convertRuleQueryResultToIssueUploadRequest(r *RuleQueryResult) *issue.UploadIssueStream_Request {

	// rule
	rule := &issue.UploadIssueStream_Rule{
		Name:   r.RuleBlock.Name,
		Query:  r.RuleBlock.Query,
		Labels: r.RuleBlock.Labels,
		Output: r.Row.String(),
	}
	if r.RuleBlock.MetadataBlock != nil {
		rule.Metadata = &issue.UploadIssueStream_Metadata{
			Id:          r.RuleBlock.MetadataBlock.Id,
			Author:      r.RuleBlock.MetadataBlock.Author,
			Description: r.RuleBlock.MetadataBlock.Description,
			Provider:    r.RuleBlock.MetadataBlock.Provider,
			Remediation: r.RuleBlock.MetadataBlock.Remediation,
			Severity:    x.ruleSeverity(r.RuleBlock.MetadataBlock.Severity),
			Tags:        r.RuleBlock.MetadataBlock.Tags,
			Title:       r.RuleBlock.MetadataBlock.Title,
		}
	}

	// provider
	ruleProvider := &issue.UploadIssueStream_Provider{
		Provider: r.Provider.Name,
		Version:  r.Provider.Version,
	}
	if r.ProviderConfiguration != nil {
		ruleProvider.Name = r.ProviderConfiguration.Name
	} else {
		ruleProvider.Name = "NOT_CONFIGURATION"
	}

	// module
	ruleModule := &issue.UploadIssueStream_Module{
		Name:             r.Module.BuildFullName(),
		Source:           r.Module.Source,
		DependenciesPath: r.Module.DependenciesPath,
	}

	// context
	ruleContext := &issue.UploadIssueStream_Context{
		SrcTableNames: r.RulePlan.BindingTables,
		Schema:        pgstorage.GetSchemaKey(ruleProvider.Provider, ruleProvider.Version, r.ProviderConfiguration),
	}

	index := x.issueIdGenerator.Add(1)
	return &issue.UploadIssueStream_Request{
		Index:    int32(index),
		Rule:     rule,
		Provider: ruleProvider,
		Module:   ruleModule,
		Context:  ruleContext,
	}
}

// Convert the original level to the enumerated value accepted by the cloud
func (x *ProjectCloudLifeCycleExecutor) ruleSeverity(severity string) issue.UploadIssueStream_Severity {
	switch strings.ToUpper(severity) {
	case "INFORMATIONAL":
		return issue.UploadIssueStream_INFORMATIONAL
	case "LOW":
		return issue.UploadIssueStream_LOW
	case "MEDIUM":
		return issue.UploadIssueStream_MEDIUM
	case "HIGH":
		return issue.UploadIssueStream_HIGH
	case "CRITICAL":
		return issue.UploadIssueStream_CRITICAL
	case "UNKNOWN":
		return issue.UploadIssueStream_UNKNOWN
	default:
		return issue.UploadIssueStream_UNKNOWN
	}
}

// ------------------------------------------------ ---------------------------------------------------------------------

// UploadLog add log to send cloud waitting queue
func (x *ProjectCloudLifeCycleExecutor) UploadLog(ctx context.Context, diagnostics *schema.Diagnostics) bool {

	if utils.IsEmpty(diagnostics) {
		return false
	}

	// show in console & log file
	x.options.MessageChannel.Send(diagnostics)

	// send to cloud
	if x.logStreamUploader == nil {
		logger.ErrorF("logStreamUploader is nil, ignore upload log")
		return utils.HasError(diagnostics)
	}
	for _, d := range diagnostics.GetDiagnosticSlice() {
		id := x.logIdGenerator.Add(1)
		isSubmitSuccess, d := x.logStreamUploader.Submit(ctx, int(id), &log.UploadLogStream_Request{
			Index: uint64(id),
			Stage: x.stage,
			Msg:   d.Content(),
			Level: x.toGrpcLevel(d.Level()),
			Time:  timestamppb.Now(),
		})
		x.options.MessageChannel.Send(d)
		if !isSubmitSuccess {
			logger.ErrorF("submit log index %d to uploader failed", id)
		} else {
			logger.InfoF("submit log index %d to uploader success", id)
		}
	}
	return utils.HasError(diagnostics)
}

func (x *ProjectCloudLifeCycleExecutor) toGrpcLevel(level schema.DiagnosticLevel) log.Level {
	switch level {
	case schema.DiagnosisLevelTrace:
		return log.Level_LEVEL_DEBUG
	case schema.DiagnosisLevelDebug:
		return log.Level_LEVEL_DEBUG
	case schema.DiagnosisLevelInfo:
		return log.Level_LEVEL_INFO
	case schema.DiagnosisLevelWarn:
		return log.Level_LEVEL_WARN
	case schema.DiagnosisLevelError:
		return log.Level_LEVEL_ERROR
	case schema.DiagnosisLevelFatal:
		return log.Level_LEVEL_FATAL
	default:
		return log.Level_LEVEL_INFO
	}
}

// ------------------------------------------------ ---------------------------------------------------------------------

// ShutdownAndWait close send queue and wait uploader done
func (x *ProjectCloudLifeCycleExecutor) ShutdownAndWait(ctx context.Context) {

	if x.logStreamUploader != nil {
		x.logStreamUploader.ShutdownAndWait(ctx)
		x.logStreamUploader.GetOptions().MessageChannel.ReceiverWait()
	}

	if x.issueStreamUploader != nil {
		x.issueStreamUploader.ShutdownAndWait(ctx)
		x.issueStreamUploader.GetOptions().MessageChannel.ReceiverWait()
	}

	x.options.MessageChannel.SenderWaitAndClose()

}

func (x *ProjectCloudLifeCycleExecutor) ChangeLogStage(stage log.StageType) {
	// change self first
	x.stage = stage
}

// ReportTaskStatus Modify the current state of the task
func (x *ProjectCloudLifeCycleExecutor) ReportTaskStatus(stage log.StageType, status log.Status) {

	if x.logClient == nil {
		logger.ErrorF("can not change task log status, not login")
		return
	}
	logger.InfoF("begin change task log status, stage = %d, status = %d", stage, status)
	logStatus, err := x.logClient.UploadLogStatus(x.cloudClient.BuildMetaContext(), &log.UploadLogStatus_Request{
		Stage:  stage,
		Status: status,
		Time:   timestamppb.Now(),
	})
	if err != nil {
		logger.ErrorF("change task log status error: %s, stage = %d, status = %d", err.Error(), stage, status)
		return
	}
	if logStatus.Diagnosis != nil && logStatus.Diagnosis.Code != 0 {
		logger.ErrorF("change task log status response error, code = %d, message = %s", logStatus.Diagnosis.Code, logStatus.Diagnosis.Msg)
	} else {
		logger.InfoF("change task log status success, stage = %d, status = %d", stage, status)
	}
}

// ------------------------------------------------- --------------------------------------------------------------------
