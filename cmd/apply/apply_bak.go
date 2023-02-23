package apply

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/selefra/selefra-provider-sdk/provider/schema"
//	"github.com/selefra/selefra-provider-sdk/storage"
//	"github.com/selefra/selefra/cli_ui"
//	"github.com/selefra/selefra/cmd/login"
//	"github.com/selefra/selefra/cmd/provider"
//	"github.com/selefra/selefra/cmd/tools"
//	"github.com/selefra/selefra/config"
//	"github.com/selefra/selefra/global"
//	"github.com/selefra/selefra/pkg/cloud_sdk"
//	"github.com/selefra/selefra/pkg/storage/pgstorage"
//
//	"github.com/spf13/cobra"
//	yaml "gopkg.in/yaml.v3"
//	"io"
//	"os"
//	"path"
//	"path/filepath"
//	"strings"
//	"text/template"
//	"time"
//)
//
//func NewApplyCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:              "apply",
//		Short:            "Analyze infrastructure",
//		Long:             "Analyze infrastructure",
//		PersistentPreRun: global.DefaultWrappedInit(),
//		RunE:             apply,
//	}
//
//	cmd.SetHelpFunc(cmd.HelpFunc())
//	return cmd
//}
//
//func apply(cmd *cobra.Command, args []string) error {
//	ctx := cmd.Context()
//	_ = login.ShouldLogin()
//
//	rootConfig, err := config.GetConfig()
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return err
//	}
//
//	// relvPrjName is the cloud relevant project name for current project
//	relvPrjName := global.RelvPrjName()
//
//	_, err = http_client.TryCreateProject(relvPrjName)
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return nil
//	}
//	taskRes, err := http_client.TryCreateTask(relvPrjName)
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return nil
//	}
//	if taskRes != nil {
//		cloud_sdk.SetTaskID(taskRes.Data.TaskUUID)
//	}
//
//	global.SetStage("initializing")
//
//	lockArr, err := provider.Sync(ctx)
//	defer func() {
//		for _, item := range lockArr {
//			err := item.Storage.UnLock(context.Background(), item.SchemaKey, item.Uuid)
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//			}
//		}
//	}()
//	if err != nil {
//		_ = http_client.TrySetUpStage(relvPrjName, http_client.Failed)
//		cli_ui.Errorln(err.Error())
//		return nil
//	}
//
//	var project string
//	if relvPrjName != "" {
//		project = relvPrjName
//	} else {
//		project = ""
//	}
//
//	if _, err := cloud_sdk.UploadLogStatus(); err != nil {
//		cli_ui.Errorln(err)
//	}
//
//	global.SetStage("infrastructure")
//
//	for _, decl := range rootConfig.Selefra.ProviderDecls {
//		prvds := tools.ProvidersByID(rootConfig, decl.Name)
//		for _, prvd := range prvds {
//			schemaKey := config.GetSchemaKey(decl, *prvd)
//			storage, diag := pgstorage.Storage(ctx, pgstorage.WithSearchPath(schemaKey))
//			if diag != nil {
//				err := cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
//				if err != nil {
//					return fmt.Errorf("failed to create pgstorage")
//				}
//			}
//
//			modules, err := config.GetModules()
//			if err != nil {
//				err = http_client.TrySetUpStage(relvPrjName, http_client.Failed)
//				cli_ui.Errorln("Client creation error:" + err.Error())
//				return nil
//			}
//			cli_ui.Successln(`----------------------------------------------------------------------------------------------
//
//Loading Selefra analysis code ...`)
//
//			var mRules []config.Rule
//			if len(modules) == 0 {
//				mRules = GetAllRules()
//			} else {
//				mRules = GetRules(modules)
//			}
//
//			cli_ui.Successf("\n---------------------------------- Result for rules  ----------------------------------------\n")
//
//			err = RunRules(ctx, rootConfig, storage, project, mRules, schemaKey)
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				return nil
//			}
//
//		}
//	}
//
//	if _, err := cloud_sdk.UploadLogStatus(); err != nil {
//		cli_ui.Errorln(err)
//	}
//
//	err = UploadWorkspace(relvPrjName)
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		sErr := http_client.TrySetUpStage(relvPrjName, http_client.Failed)
//		if sErr != nil {
//			cli_ui.Errorln(sErr.Error())
//		}
//		return nil
//	}
//	return nil
//}
//
//func UploadWorkspace(project string) error {
//	fileMap, err := config.FileMap(global.WorkSpace())
//	if err != nil {
//		return err
//	}
//	err = http_client.TryUploadWorkspace(project, fileMap)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func getTableMap(tableMap map[string]bool, schemaTable []*schema.Table) {
//	for i := range schemaTable {
//		tableMap[schemaTable[i].TableName] = true
//		if len(schemaTable[i].SubTables) > 0 {
//			getTableMap(tableMap, schemaTable[i].SubTables)
//		}
//	}
//}
//
//func match(s string, whitelistWordSet map[string]bool) []string {
//	var matchResultSet []string
//	inWord := false
//	lastIndex := 0
//	for index, c := range s {
//		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9' {
//			if !inWord {
//				inWord = true
//				lastIndex = index
//			}
//		} else {
//			if inWord {
//				word := s[lastIndex:index]
//				if _, exists := whitelistWordSet[word]; exists {
//					matchResultSet = append(matchResultSet, word)
//				}
//				inWord = false
//			}
//		}
//	}
//	return matchResultSet
//}
//
//func getSqlTables(sql string, tableMap map[string]bool) (tables []string) {
//	nonStr := strings.Replace(sql, "\n", "", -1)
//	tables = match(nonStr, tableMap)
//	return tables
//}
//
//func UploadIssueFunc(ctx context.Context, IssueReq <-chan *issue.Req, ticker *time.Ticker) {
//	for {
//		if ticker != nil {
//			ticker.Reset(30 * time.Second)
//		}
//		select {
//		case req := <-IssueReq:
//			if err := cloud_sdk.IssueStreamSend(req); err != nil {
//				cli_ui.Errorf("send issue to server error: %s", err.Error())
//				return
//			}
//		case <-ctx.Done():
//			_ = cloud_sdk.IssueStreamClose()
//			cli_ui.Infoln("End of reporting issue")
//			return
//		}
//	}
//}
//
//func RunRules(ctx context.Context, rootConfig *config.RootConfig, storage storage.Storage, project string, rules []config.Rule, schema string) error {
//	issueCtx, issueCancel := context.WithCancel(context.Background())
//	defer issueCancel()
//	issueChan := make(chan *issue.Req, 100)
//	ticker := time.NewTicker(30 * time.Second)
//	defer ticker.Stop()
//
//	go func() {
//		select {
//		case <-ticker.C:
//			cli_ui.Errorln("Report issue timeout")
//			_, _ = cloud_sdk.UploadLogStatus()
//			issueCancel()
//			return
//		case <-issueCtx.Done():
//			return
//		}
//	}()
//
//	go UploadIssueFunc(issueCtx, issueChan, ticker)
//
//	for _, rule := range rules {
//		var variablesMap = make(map[string]interface{})
//		for i := range rootConfig.Variables {
//			variablesMap[rootConfig.Variables[i].Key] = rootConfig.Variables[i].Default
//		}
//		queryStr, err := fmtTemplate(rule.Query, variablesMap)
//		res, diag := storage.Query(ctx, queryStr)
//		if diag != nil {
//			cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
//			continue
//		}
//		table, diag := res.ReadRows(-1)
//		if diag != nil {
//			cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
//			continue
//		}
//		column := table.GetColumnNames()
//		rows := table.GetMatrix()
//		if len(rows) == 0 {
//			continue
//		}
//		cli_ui.Successf("%rootConfig - Rule \"%rootConfig\"\n", rule.Path, rule.Name)
//		cli_ui.Successln("Schema:")
//		cli_ui.Successln(schema + "\n")
//		cli_ui.Successln("Description:")
//
//		desc, err := fmtTemplate(rule.Metadata.Description, variablesMap)
//		if err != nil {
//			cli_ui.Errorln(err.Error())
//			return err
//		}
//		cli_ui.Successln("	" + desc)
//
//		cli_ui.Successln("Policy:")
//		schemaTables, schemaDiag := storage.TableList(ctx, schema)
//		if schemaDiag != nil {
//			err := cli_ui.PrintDiagnostic(schemaDiag.GetDiagnosticSlice())
//			if err != nil {
//				return err
//			}
//		}
//		var tableMap = make(map[string]bool)
//		getTableMap(tableMap, schemaTables)
//
//		uploadTables := getSqlTables(queryStr, tableMap)
//		if err != nil {
//			cli_ui.Errorln(err.Error())
//			return err
//		}
//		cli_ui.Successln("	" + queryStr)
//
//		cli_ui.Successln("Output")
//		for _, row := range rows {
//			var outMetaData issue.Metadata
//			var baseRow = make(map[string]interface{})
//			var outPut = rule.Output
//			var outMap = make(map[string]interface{})
//			for index, value := range row {
//				key := column[index]
//				outMap[key] = value
//			}
//			baseRow = outMap
//			baseRowStr, err := json.Marshal(baseRow)
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				return err
//			}
//			var outByte bytes.Buffer
//			err = json.Indent(&outByte, baseRowStr, "", "\t")
//			out, err := fmtTemplate(outPut, outMap)
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				return err
//			}
//			var remediation string
//			var remediationPath string
//			if filepath.IsAbs(rule.Metadata.Remediation) {
//				remediationPath = rule.Metadata.Remediation
//			} else {
//				remediationPath = filepath.Join(global.WorkSpace(), rule.Metadata.Remediation)
//			}
//			remediationByte, err := os.ReadFile(remediationPath)
//			remediation = string(remediationByte)
//			if err != nil {
//				remediation = err.Error()
//			}
//			outMetaData = issue.Metadata{
//				Id:           rule.Metadata.Id,
//				Severity:     rule.Metadata.Severity,
//				Remediation:  remediation,
//				Tags:         rule.Metadata.Tags,
//				SrcTableName: uploadTables,
//				Provider:     rule.Metadata.Provider,
//				Title:        rule.Metadata.Title,
//				Author:       rule.Metadata.Author,
//				Description:  desc,
//				Output:       outByte.String(),
//			}
//
//			cli_ui.Successln("	" + out)
//
//			var outLabel = make(map[string]string)
//			for key := range rule.Labels {
//				switch rule.Labels[key].(type) {
//				case string:
//					outStr, _ := fmtTemplate(rule.Labels[key].(string), baseRow)
//					outLabel[key] = outStr
//				case []string:
//					var out []string
//					for _, v := range rule.Labels[key].([]string) {
//						s, _ := fmtTemplate(v, baseRow)
//						out = append(out, s)
//					}
//					outLabel[key] = strings.Join(out, ",")
//				}
//			}
//
//			reqs := issue.Req{
//				Name:        rule.Name,
//				Query:       rule.Query,
//				Metadata:    &outMetaData,
//				Labels:      outLabel,
//				ProjectName: project,
//				TaskUUID:    cloud_sdk.TaskID(),
//				Token:       cloud_sdk.Token(),
//				Schema:      schema,
//			}
//			select {
//			case <-issueCtx.Done():
//				return nil
//			default:
//				issueChan <- &reqs
//			}
//		}
//	}
//	return nil
//}
//
//// GetAllRules get all rules from workspace
//func GetAllRules() []config.Rule {
//	rules, _ := config.GetRules()
//	for i := range rules.Rules {
//		if strings.HasPrefix(rules.Rules[i].Query, ".") {
//			sqlByte, err := os.ReadFile(filepath.Join(".", rules.Rules[i].Query))
//			if err != nil {
//				cli_ui.Errorf("sql open error:%s", err.Error())
//				return nil
//			}
//			rules.Rules[i].Query = string(sqlByte)
//		}
//	}
//	return rules.Rules
//}
//
//// GetRules find all rules in modules
//func GetRules(modules []config.Module) []config.Rule {
//	var rules []config.Rule
//	for _, module := range modules {
//		if rule := GetModuleRules(module); rule != nil {
//			rules = append(rules, rule...)
//		}
//	}
//	return rules
//}
//
//// GetModuleRules find all rules according to given module
//func GetModuleRules(module config.Module) []config.Rule {
//	var resRule config.RuleSet
//	var b []byte
//	var err error
//
//	for _, use := range module.Uses {
//		var usePath string
//		if path.IsAbs(use) || strings.Index(use, "://") > -1 {
//			usePath = use
//		} else {
//			usePath = filepath.Join(global.WorkSpace(), use)
//		}
//		if strings.Index(usePath, "://") > -1 {
//			d := config.Downloader{Url: usePath}
//			b, err = d.Get()
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				return nil
//			}
//		} else {
//			b, err = os.ReadFile(usePath)
//			if err != nil {
//				cli_ui.Errorln(err.Error())
//				return nil
//			}
//		}
//
//		var baseRule config.RuleSet
//		err = yaml.Unmarshal(b, &baseRule)
//		if err != nil {
//			cli_ui.Errorln(err.Error())
//			return nil
//		}
//
//		if err != nil {
//			cli_ui.Errorln(err.Error())
//			return nil
//		}
//		var ruleConfig config.RuleSet
//		err = yaml.Unmarshal([]byte(string(b)), &ruleConfig)
//		if err != nil {
//			cli_ui.Errorln(err.Error())
//			return nil
//		}
//		for i := range ruleConfig.Rules {
//			ruleConfig.Rules[i].Output = baseRule.Rules[i].Output
//			ruleConfig.Rules[i].Query = baseRule.Rules[i].Query
//			ruleConfig.Rules[i].Path = use
//			_, err := os.Stat(filepath.Join(global.WorkSpace(), ruleConfig.Rules[i].Query))
//			if err == nil {
//				var sqlPath string
//				if filepath.IsAbs(ruleConfig.Rules[i].Query) {
//					sqlPath = ruleConfig.Rules[i].Query
//				} else {
//					sqlPath = filepath.Join(global.WorkSpace(), ruleConfig.Rules[i].Query)
//				}
//				sqlByte, err := os.ReadFile(sqlPath)
//				if err != nil {
//					cli_ui.Errorf("sql open error:%s", err.Error())
//					return nil
//				}
//				ruleConfig.Rules[i].Query = string(sqlByte)
//			}
//			cli_ui.Successf("	%s - Rule %s: loading ... \n", use, baseRule.Rules[i].Name)
//		}
//		resRule.Rules = append(resRule.Rules, ruleConfig.Rules...)
//	}
//	return resRule.Rules
//}
//
//func fmtTemplate(temp string, params map[string]interface{}) (string, error) {
//	t, err := template.New("temp").Parse(temp)
//	if err != nil {
//		cli_ui.Errorln("Format rule template error:", err.Error())
//		return "", err
//	}
//	b := bytes.Buffer{}
//	err = t.Execute(&b, params)
//	if err != nil {
//		cli_ui.Errorln("Format rule template error:", err.Error())
//		return "", err
//	}
//	by, err := io.ReadAll(&b)
//	if err != nil {
//		cli_ui.Errorln("Format rule template error:", err.Error())
//		return "", err
//	}
//	return string(by), nil
//}
