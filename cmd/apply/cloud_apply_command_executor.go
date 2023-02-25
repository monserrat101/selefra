package apply

//import (
//	"context"
//	"fmt"
//	"github.com/selefra/selefra-provider-sdk/env"
//	"github.com/selefra/selefra-provider-sdk/provider/schema"
//	"github.com/selefra/selefra/cli_ui"
//	"github.com/selefra/selefra/pkg/cloud_sdk"
//	selefraGrpc "github.com/selefra/selefra/pkg/grpc"
//	"github.com/selefra/selefra/pkg/grpc/pb/issue"
//	"github.com/selefra/selefra/pkg/grpc/pb/log"
//	"github.com/selefra/selefra/pkg/logger"
//	"github.com/selefra/selefra/pkg/message"
//	"github.com/selefra/selefra/pkg/modules/executors"
//	"github.com/selefra/selefra/pkg/modules/module"
//	"github.com/selefra/selefra/pkg/storage/pgstorage"
//	"github.com/selefra/selefra/pkg/utils"
//	"google.golang.org/protobuf/types/known/timestamppb"
//	"strings"
//	"sync/atomic"
//)
//
//// ------------------------------------------------- --------------------------------------------------------------------
//
//type CloudApplyCommandExecutor struct {
//
//	// The address of the Cloud cluster to which you are connecting
//	cloudHost string
//
//	// the client for connect to cloud
//	client *cloud_sdk.CloudClient
//
//	// for upload log
//	logClient         log.LogClient
//	logStreamUploader *selefraGrpc.StreamUploader[log.Log_UploadLogStreamClient, int, *log.UploadLogStream_Request, *log.UploadLogStream_Response]
//
//	// for upload issue
//	issueStreamUploader *selefraGrpc.StreamUploader[issue.Issue_UploadIssueStreamClient, int, *issue.UploadIssueStream_Request, *issue.UploadIssueStream_Response]
//
//	// index generator
//	logIdGenerator   atomic.Int64
//	issueIdGenerator atomic.Int64
//
//	// task current stage
//	stage log.StageType
//}
//
//func NewCloudApplyCommandExecutor(cloudHost string) *CloudApplyCommandExecutor {
//	return &CloudApplyCommandExecutor{
//		cloudHost: cloudHost,
//	}
//}
//
//func (x *CloudApplyCommandExecutor) getDSN(ctx context.Context, rootModule *module.Module) string {
//
//	// Use the configuration of the current module first if it is configured in the current module
//	if rootModule.SelefraBlock != nil && rootModule.SelefraBlock.ConnectionBlock != nil {
//		logger.InfoF("get dsn from selefra block")
//		return rootModule.SelefraBlock.ConnectionBlock.BuildDSN()
//	}
//
//	// If you log in, use the remote configuration
//	if x.client.IsLoggedIn() {
//		dsn, diagnostics := x.client.FetchOrgDSN()
//		if err := x.UploadLog(ctx, diagnostics); err == nil && dsn != "" {
//			logger.InfoF("get dsn from cloud")
//			return dsn
//		}
//	}
//
//	// Environment variable
//	if env.GetDatabaseDsn() != "" {
//		logger.InfoF("get dsn from env")
//		return env.GetDatabaseDsn()
//	}
//
//	// Use the built-in PG database
//	logger.InfoF("get dsn use default")
//	// TODO
//	return ""
//}
//
//func (x *CloudApplyCommandExecutor) initCloudClient(ctx context.Context, module *module.Module) bool {
//
//	// create cloud client
//	client, d := cloud_sdk.NewCloudClient(x.cloudHost)
//	if err := cli_ui.PrintDiagnostics(d); err != nil {
//		return false
//	}
//	x.client = client
//
//	// find local cloud token
//	credentials, _ := client.GetCredentials()
//	if credentials != nil {
//		login, diagnostics := client.Login(credentials.Token)
//		if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
//			cli_ui.ShowLoginFailed(credentials.Token)
//			return false
//		} else {
//
//			// login success
//			cli_ui.ShowLoginSuccess(x.cloudHost, login)
//
//			// relative project
//			if module.SelefraBlock == nil || module.SelefraBlock.CloudBlock == nil || module.SelefraBlock.CloudBlock.Project == "" {
//				cli_ui.Errorln("Failed to connect to the cloud, you must specify the project name")
//				return false
//			}
//			project, d := client.CreateProject(module.SelefraBlock.CloudBlock.Project)
//			if err := cli_ui.PrintDiagnostics(d); err != nil {
//				return false
//			}
//			cli_ui.Successf("Successfully connected to cloud, associated to project %s\n", project.Name)
//
//			// create task
//			task, d := client.CreateTask(project.Name)
//			if err := cli_ui.PrintDiagnostics(d); err != nil {
//				cli_ui.Errorf("Failed to create a task for the project %s\n", project.Name)
//				return false
//			}
//			cli_ui.Errorf("Succeeded in creating a task, id = %s\n", task.TaskId)
//
//			x.initLogUploader(client)
//			x.initIssueUploader(client)
//
//			// change task status to begin
//			x.ChangeTaskLogStatus(log.StageType_STAGE_TYPE_INITIALIZING, log.Status_STATUS_SUCCESS)
//
//			// show log
//			_ = x.UploadLog(ctx, schema.NewDiagnostics().AddInfo("begin init task id %s", task.TaskId))
//		}
//	}
//
//	return true
//}
//
//// init issue uploader for send issue to cloud
//func (x *CloudApplyCommandExecutor) initIssueUploader(client *cloud_sdk.CloudClient) bool {
//	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
//		_ = cli_ui.PrintDiagnostics(message)
//	})
//	issueStreamUploader, diagnostics := client.NewIssueStreamUploader(messageChannel)
//	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
//		return false
//	}
//	issueStreamUploader.RunUploaderWorker()
//	x.issueStreamUploader = issueStreamUploader
//	return true
//}
//
//// init log uploader for send log to loud
//func (x *CloudApplyCommandExecutor) initLogUploader(client *cloud_sdk.CloudClient) bool {
//	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
//		_ = cli_ui.PrintDiagnostics(message)
//	})
//	logClient, logStreamUploader, diagnostics := client.NewLogStreamUploader(messageChannel)
//	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
//		return false
//	}
//	logStreamUploader.RunUploaderWorker()
//	x.logClient = logClient
//	x.logStreamUploader = logStreamUploader
//	return true
//}
//
//// UploadIssue add issue to send cloud queue
//func (x *CloudApplyCommandExecutor) UploadIssue(ctx context.Context, r *executors.RuleQueryResult) {
//
//	// TODO Modified for unified display
//	//send to console & log file
//	//var outByte bytes.Buffer
//	//err := json.Indent(&outByte, json_util.ToJsonBytes(r.RuleBlock), "", "\t")
//	//if err != nil {
//	//	logger.ErrorF("format issue error: %s", err.Error())
//	//} else {
//	//	cli_ui.Successln(outByte.String())
//	//}
//	var consoleOutput strings.Builder
//	consoleOutput.WriteString(fmt.Sprintf("Rule name %s, ", r.RuleBlock.Name))
//	if r.RuleBlock.MetadataBlock != nil && r.RuleBlock.MetadataBlock.Id != "" {
//		consoleOutput.WriteString(fmt.Sprintf("id %s, ", r.RuleBlock.MetadataBlock.Id))
//	}
//	consoleOutput.WriteString(fmt.Sprintf("output %s", r.RuleBlock.Output))
//	cli_ui.Successln(consoleOutput.String())
//
//	// send to cloud
//	if x.issueStreamUploader == nil {
//		logger.ErrorF("issueStreamUploader is nil")
//		return
//	}
//	request := x.convertRuleQueryResultToIssueUploadRequest(r)
//	x.issueStreamUploader.Submit(ctx, int(request.Index), request)
//}
//
//// Convert the query results of the rules into a format uploaded to the Cloud
//func (x *CloudApplyCommandExecutor) convertRuleQueryResultToIssueUploadRequest(r *executors.RuleQueryResult) *issue.UploadIssueStream_Request {
//
//	// rule
//	rule := &issue.UploadIssueStream_Rule{
//		Name:   r.RuleBlock.Name,
//		Query:  r.RuleBlock.Query,
//		Labels: r.RuleBlock.Labels,
//		Output: r.Row.String(),
//	}
//	if r.RuleBlock.MetadataBlock != nil {
//		rule.Metadata = &issue.UploadIssueStream_Metadata{
//			Id:          r.RuleBlock.MetadataBlock.Id,
//			Author:      r.RuleBlock.MetadataBlock.Author,
//			Description: r.RuleBlock.MetadataBlock.Description,
//			Provider:    r.RuleBlock.MetadataBlock.Provider,
//			Remediation: r.RuleBlock.MetadataBlock.Remediation,
//			Severity:    x.ruleSeverity(r.RuleBlock.MetadataBlock.Severity),
//			Tags:        r.RuleBlock.MetadataBlock.Tags,
//			Title:       r.RuleBlock.MetadataBlock.Title,
//		}
//	}
//
//	// provider
//	ruleProvider := &issue.UploadIssueStream_Provider{
//		Provider: r.Provider.Name,
//		Version:  r.Provider.Version,
//		//Name:     "",
//	}
//	if r.ProviderConfiguration != nil {
//		ruleProvider.Name = r.ProviderConfiguration.Name
//	}
//
//	// module
//	ruleModule := &issue.UploadIssueStream_Module{
//		Name:             r.Module.BuildFullName(),
//		Source:           r.Module.Source,
//		DependenciesPath: r.Module.DependenciesPath,
//	}
//
//	// context
//	ruleContext := &issue.UploadIssueStream_Context{
//		SrcTableNames: r.RulePlan.BindingTables,
//		Schema:        pgstorage.GetSchemaKey(ruleProvider.Provider, ruleProvider.Version, r.ProviderConfiguration),
//	}
//
//	return &issue.UploadIssueStream_Request{
//		Index:    int32(r.Index),
//		Rule:     rule,
//		Provider: ruleProvider,
//		Module:   ruleModule,
//		Context:  ruleContext,
//	}
//}
//
//// Convert the original level to the enumerated value accepted by the cloud
//func (x *CloudApplyCommandExecutor) ruleSeverity(severity string) issue.UploadIssueStream_Severity {
//	switch strings.ToUpper(severity) {
//	case "INFORMATIONAL":
//		return issue.UploadIssueStream_INFORMATIONAL
//	case "LOW":
//		return issue.UploadIssueStream_LOW
//	case "MEDIUM":
//		return issue.UploadIssueStream_MEDIUM
//	case "HIGH":
//		return issue.UploadIssueStream_HIGH
//	case "CRITICAL":
//		return issue.UploadIssueStream_CRITICAL
//	case "UNKNOWN":
//		return issue.UploadIssueStream_UNKNOWN
//	default:
//		return issue.UploadIssueStream_UNKNOWN
//	}
//}
//
//// ------------------------------------------------ ---------------------------------------------------------------------
//
//// UploadLog add log to send cloud waitting queue
//func (x *CloudApplyCommandExecutor) UploadLog(ctx context.Context, diagnostics *schema.Diagnostics) error {
//
//	if utils.IsEmpty(diagnostics) {
//		return nil
//	}
//
//	// show in console & log file
//	err := cli_ui.PrintDiagnostics(diagnostics)
//
//	// send to cloud
//	if x.logStreamUploader == nil {
//		logger.ErrorF("logStreamUploader is nil")
//		return err
//	}
//	for _, d := range diagnostics.GetDiagnosticSlice() {
//		id := x.logIdGenerator.Add(1)
//		isSubmitSuccess, d := x.logStreamUploader.Submit(ctx, int(id), &log.UploadLogStream_Request{
//			Index: uint64(id),
//			Stage: x.stage,
//			Msg:   d.Content(),
//			Level: x.toGrpcLevel(d.Level()),
//			Time:  timestamppb.Now(),
//		})
//		_ = cli_ui.PrintDiagnostics(d)
//		if !isSubmitSuccess {
//			logger.ErrorF("submit log index %d to uploader failed", id)
//		}
//	}
//	return err
//}
//
//func (x *CloudApplyCommandExecutor) toGrpcLevel(level schema.DiagnosticLevel) log.Level {
//	switch level {
//	case schema.DiagnosisLevelTrace:
//		return log.Level_LEVEL_DEBUG
//	case schema.DiagnosisLevelDebug:
//		return log.Level_LEVEL_DEBUG
//	case schema.DiagnosisLevelInfo:
//		return log.Level_LEVEL_INFO
//	case schema.DiagnosisLevelWarn:
//		return log.Level_LEVEL_WARN
//	case schema.DiagnosisLevelError:
//		return log.Level_LEVEL_ERROR
//	case schema.DiagnosisLevelFatal:
//		return log.Level_LEVEL_FATAL
//	default:
//		return log.Level_LEVEL_INFO
//	}
//}
//
//// ------------------------------------------------ ---------------------------------------------------------------------
//
//// ShutdownAndWait close send queue and wait uploader done
//func (x *CloudApplyCommandExecutor) ShutdownAndWait(ctx context.Context) {
//
//	if x.logStreamUploader != nil {
//		x.logStreamUploader.ShutdownAndWait(ctx)
//		x.logStreamUploader.GetOptions().MessageChannel.ReceiverWait()
//	}
//
//	if x.issueStreamUploader != nil {
//		x.issueStreamUploader.ShutdownAndWait(ctx)
//		x.issueStreamUploader.GetOptions().MessageChannel.ReceiverWait()
//	}
//
//}
//
//// ------------------------------------------------- --------------------------------------------------------------------
//
//// ChangeTaskLogStatus Modify the current state of the task
//func (x *CloudApplyCommandExecutor) ChangeTaskLogStatus(stage log.StageType, status log.Status) {
//	if x.logClient == nil {
//		logger.ErrorF("can not change task log status, not login")
//		return
//	}
//	logStatus, err := x.logClient.UploadLogStatus(x.client.BuildMetaContext(), &log.UploadLogStatus_Request{
//		Stage:  stage,
//		Status: status,
//		Time:   timestamppb.Now(),
//	})
//	if err != nil {
//		logger.ErrorF("change task log status error: %s, stage = %d, status = %d", err.Error(), stage, status)
//		return
//	}
//	if logStatus.Diagnosis != nil && logStatus.Diagnosis.Code != 0 {
//		logger.ErrorF("change task log status response error, code = %d, message = %s", logStatus.Diagnosis.Code, logStatus.Diagnosis.Msg)
//	} else {
//		logger.InfoF("change task log status success, stage = %d, status = %d", stage, status)
//	}
//}
//
//// ------------------------------------------------- --------------------------------------------------------------------
