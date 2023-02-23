package init

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/cmd/version"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/cloud_sdk"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/parser"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/storage/pgstorage"
	"github.com/selefra/selefra/pkg/utils"
	version2 "github.com/selefra/selefra/pkg/version"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InitCommandExecutorOptions struct {
	DownloadWorkspace string
	ProjectWorkspace  string
	IsForceInit       bool
	RelevanceProject  string
	DSN               string
}

type InitCommandExecutor struct {
	cloudClient *cloud_sdk.CloudClient

	options *InitCommandExecutorOptions
}

func NewInitCommandExecutor(options *InitCommandExecutorOptions) *InitCommandExecutor {
	return &InitCommandExecutor{
		options: options,
	}
}

func (x *InitCommandExecutor) Run(ctx context.Context) {

	if !x.checkWorkspace() {
		return
	}

	// init files
	selefraBlock := x.initSelefraYaml()
	if selefraBlock != nil {
		x.initProvidersYaml(ctx, selefraBlock.RequireProvidersBlock)
	}
	x.initRulesYaml()
	x.initModulesYaml()

	cli_ui.Successf("Initializing workspace done.\n")
}

func (x *InitCommandExecutor) checkWorkspace() bool {

	// 1. check if workspace dir exist
	_, err := os.Stat(x.options.ProjectWorkspace)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(x.options.ProjectWorkspace, 0755)
		if err != nil {
			cli_ui.Errorf("create workspace directory: %s failed: %s", x.options.ProjectWorkspace, err.Error())
			return false
		}
	}
	dir, _ := os.ReadDir(x.options.ProjectWorkspace)
	for i, v := range dir { // ignore logs dir
		if v.Name() == "logs" {
			dir = append(dir[0:i], dir[i+1:]...)
		}
	}

	// 2. workspace must be empty or set force flag
	//force, _ := cmd.PersistentFlags().GetBool("force")
	if len(dir) != 0 {
		if !x.options.IsForceInit {
			cli_ui.Errorf("%s is not empty; Rerun in an empty directory, or use -- force/-f to force overwriting in the current directory\n", x.options.ProjectWorkspace)
			return false
		} else if !x.reInit() {
			return false
		}
	}

	return true
}

const (
	SelefraInputInitForceConfirm     = "SELEFRA_INPUT_INIT_FORCE_CONFIRM"
	SelefraInputInitRelevanceProject = "SELEFRA_INPUT_INIT_RELEVANCE_PROJECT"
)

// reInit check if current workspace is selefra workspace, then tell user to choose if rewrite selefra workspace
func (x *InitCommandExecutor) reInit() bool {
	//_, err := config.GetConfig()
	//if err != nil && errors.Is(err, config.ErrNotSelefra) {
	//	return nil
	//}

	reader := bufio.NewReader(os.Stdin)
	cli_ui.Warningf("Warning: %s is already init. Continue and overwrite it?[Y/N]", x.options.ProjectWorkspace)
	text, err := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))
	if err != nil && !errors.Is(err, io.EOF) {
		cli_ui.Errorf("read you input error: %s", err.Error())
		return false
	}

	// for test
	if text == "" {
		text = os.Getenv(SelefraInputInitForceConfirm)
	}

	if text != "y" && text != "Y" {
		cli_ui.Errorf("config file already exists")
		return false
	}

	return true
}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *InitCommandExecutor) initSelefraYaml() *module.SelefraBlock {

	selefraBlock := module.NewSelefraBlock()
	projectName, b := x.getProjectName()
	if !b {
		return nil
	}
	selefraBlock.Name = projectName

	// cloud block
	selefraBlock.CloudBlock = x.getCloudBlock()

	// cli version
	selefraBlock.CliVersion = version.Version
	selefraBlock.LogLevel = "info"

	list, _ := x.chooseProvidersList()
	if len(list) > 0 {
		requiredProviderSlice := make([]*module.RequireProviderBlock, len(list))
		for index, providerName := range list {
			requiredProviderBlock := module.NewRequireProviderBlock()
			requiredProviderBlock.Name = providerName
			requiredProviderBlock.Source = providerName
			requiredProviderBlock.Version = version2.VersionLatest
			requiredProviderSlice[index] = requiredProviderBlock
		}
		selefraBlock.RequireProvidersBlock = requiredProviderSlice
	}

	selefraBlock.ConnectionBlock = x.GetConnectionBlock()

	out, err := yaml.Marshal(selefraBlock)
	if err != nil {
		cli_ui.Errorf("selefra block yaml.Marshal error: %s", err.Error())
		return nil
	}
	var selefraNode yaml.Node
	err = yaml.Unmarshal(out, &selefraNode)
	if err != nil {
		cli_ui.Errorf("selefra yaml.Unmarshal error: %s", err.Error())
		return nil
	}
	documentRoot := yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			&yaml.Node{Kind: yaml.ScalarNode, Value: parser.SelefraBlockFieldName},
			&yaml.Node{Kind: yaml.MappingNode, Content: selefraNode.Content[0].Content},
		},
	}
	marshal, err := yaml.Marshal(&documentRoot)
	if err != nil {
		cli_ui.Errorf("selefra yaml.Marshal error: %s", err.Error())
		return nil
	}
	selefraFullPath := filepath.Join(utils.AbsPath(x.options.ProjectWorkspace), "selefra.yaml")
	err = os.WriteFile(selefraFullPath, marshal, 0644)
	if err != nil {
		cli_ui.Errorf("Write %s error: %s\n", selefraFullPath, err.Error())
	} else {
		cli_ui.Successf("Write %s success\n", selefraFullPath)
	}

	return selefraBlock
}

func (x *InitCommandExecutor) getCloudBlock() *module.CloudBlock {

	cloudBlock := module.NewCloudBlock()

	// project name
	projectName, b := x.getProjectName()

	if !b {
		return nil
	}
	cloudBlock.Project = projectName

	if x.cloudClient != nil {
		credentials, diagnostics := x.cloudClient.GetCredentials()
		if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
			return nil
		}
		cloudBlock.Organization = credentials.OrgName
		cloudBlock.HostName = credentials.ServerHost
	}

	return cloudBlock
}

// init module.yaml
func (x *InitCommandExecutor) initModulesYaml() {
	const moduleComment = `
modules:
  - name: AWS_Security_Demo
    uses:
    - ./rules/
`
	moduleFullPath := filepath.Join(utils.AbsPath(x.options.ProjectWorkspace), "module.yaml")
	err := os.WriteFile(moduleFullPath, []byte(moduleComment), 0644)
	if err != nil {
		cli_ui.Errorf("Write %s error: %s\n", moduleFullPath, err.Error())
	} else {
		cli_ui.Successf("Write %s success\n", moduleFullPath)
	}
}

func (x *InitCommandExecutor) initRulesYaml() {
	const ruleComment = `
rules:
  - name: example_rule_name
    query: |
      SELECT 
        *
      FROM 
        aws_ec2_ebs_volumes 
      WHERE 
        encrypted = FALSE;
    labels:  
      resource_type: EC2 
      resource_account_id : '{{.account_id}}'
      resource_id: '{{.id}}'
      resource_region: '{{.availability_zone}}'
    metadata: 
      id: SF010302
      severity: Low
      provider: AWS
      tags:
        - Misconfigure
      author: Selefra
      remediation: remediation/ec2/ebs_volume_are_unencrypted.md
      title: EBS volume are unencrypted 
      description: Ensure that EBS volumes are encrypted.
    output: 'EBS volume are unencrypted, EBS id: {{.id}}, availability zone: {{.availability_zone}}'
`
	ruleDirectory := filepath.Join(utils.AbsPath(x.options.ProjectWorkspace), "rules")
	_ = utils.EnsureDirectoryExists(ruleDirectory)
	ruleFullPath := filepath.Join(ruleDirectory, "rule.yaml")
	err := os.WriteFile(ruleFullPath, []byte(ruleComment), 0644)
	if err != nil {
		cli_ui.Errorf("Write %s error: %s\n", ruleFullPath, err.Error())
	} else {
		cli_ui.Successf("Write %s success\n", ruleFullPath)
	}
}

func (x *InitCommandExecutor) initProvidersYaml(ctx context.Context, requiredProviders module.RequireProvidersBlock) {
	if len(requiredProviders) == 0 {
		cli_ui.Infof("No required provider, do not init providers file\n")
		return
	}
	providers, b := x.makeProviders(ctx, requiredProviders)
	if !b {
		return
	}
	out, err := yaml.Marshal(providers)
	if err != nil {
		cli_ui.Errorf("providers block yaml.Marshal error: %s", err.Error())
		return
	}
	var providersNode yaml.Node
	err = yaml.Unmarshal(out, &providersNode)
	if err != nil {
		cli_ui.Errorf("providers yaml.Unmarshal error: %s", err.Error())
		return
	}
	documentRoot := yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			&yaml.Node{Kind: yaml.ScalarNode, Value: parser.ProvidersBlockName},
			&yaml.Node{Kind: yaml.MappingNode, Content: providersNode.Content[0].Content},
		},
	}
	marshal, err := yaml.Marshal(documentRoot)
	if err != nil {
		cli_ui.Errorf("providers yaml.Marshal error: %s", err.Error())
		return
	}
	ruleFullPath := filepath.Join(utils.AbsPath(x.options.ProjectWorkspace), "providers.yaml")
	err = os.WriteFile(ruleFullPath, marshal, 0644)
	if err != nil {
		cli_ui.Errorf("Write %s error: %s\n", ruleFullPath, err.Error())
	} else {
		cli_ui.Successf("Write %s success\n", ruleFullPath)
	}
}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *InitCommandExecutor) GetConnectionBlock() *module.ConnectionBlock {

	//// 1. Try to get the DSN from the cloud
	//if x.cloudClient != nil {
	//	dsn, diagnostics := x.cloudClient.FetchOrgDSN()
	//	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
	//		return nil
	//	}
	//	if dsn != "" {
	//		return x.parseDsnAsConnectionBlock(dsn)
	//	}
	//}
	//
	//// 2.

	cli_runtime.Init(x.options.ProjectWorkspace)

	dsn, diagnostics := cli_runtime.GetDSN()
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return nil
	}
	if dsn != "" {
		return module.ParseConnectionBlockFromDSN(dsn)
	}

	return nil
}

func (x *InitCommandExecutor) getProjectName() (string, bool) {

	// 1. Use the specified one, if any
	if x.options.RelevanceProject != "" {
		return x.options.RelevanceProject, true
	}

	defaultProjectName := filepath.Base(x.options.ProjectWorkspace)

	// 2. Let the user specify from standard input, the default project name is the name of the current folder
	var err error
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("project name:(%s)", defaultProjectName)
	projectName, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		cli_ui.Errorf("read you project name error: %s", err.Error())
		return "", false
	}
	projectName = strings.TrimSpace(strings.Replace(projectName, "\n", "", -1))
	if projectName == "" {
		projectName = defaultProjectName
	}
	return projectName, true
}

func (x *InitCommandExecutor) chooseProvidersList() ([]string, bool) {
	providerNameSlice, ok := x.requestProvidersList()
	if !ok {
		return nil, false
	}
	providersSet := cli_ui.SelectProviders(providerNameSlice)
	providerSlice := make([]string, 0)
	for providerName := range providersSet {
		providerSlice = append(providerSlice, providerName)
	}
	return providerNameSlice, true
}

func (x *InitCommandExecutor) requestProvidersList() ([]string, bool) {
	var prov []string
	cli_ui.Infoln("Getting provider list...")
	req, _ := http.NewRequest("GET", "https://github.com/selefra/registry/file-list/main/provider", nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		cli_ui.Errorf("Error: %s", err.Error())
		return nil, false
	}
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		cli_ui.Errorf("Error: %s", err.Error())
		return nil, false
	}
	d.Find(".js-navigation-open.Link--primary").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "template" {
			prov = append(prov, s.Text())
		}
	})
	return prov, false
}

// ------------------------------------------------- --------------------------------------------------------------------

func (x *InitCommandExecutor) makeProviders(ctx context.Context, requiredProvidersBlock module.RequireProvidersBlock) (module.ProvidersBlock, bool) {

	providersBlock := make(module.ProvidersBlock, 0)
	// convert required provider block to
	for _, requiredProvider := range requiredProvidersBlock {

		providerInstallPlan := &planner.ProviderInstallPlan{
			Provider: registry.NewProvider(requiredProvider.Name, requiredProvider.Version),
		}

		// install providers
		messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
			_ = cli_ui.PrintDiagnostics(message)
		})
		executor, d := executors.NewProviderInstallExecutor(&executors.ProviderInstallExecutorOptions{
			Plans: []*planner.ProviderInstallPlan{
				providerInstallPlan,
			},
			MessageChannel:    messageChannel,
			DownloadWorkspace: x.options.DownloadWorkspace,
		})
		if err := cli_ui.PrintDiagnostics(d); err != nil {
			return nil, false
		}
		d = executor.Execute(ctx)
		messageChannel.ReceiverWait()
		if err := cli_ui.PrintDiagnostics(d); err != nil {
			return nil, false
		}

		// init
		configuration, b := x.getProviderInitConfiguration(ctx, executor.GetLocalProviderManager(), providerInstallPlan)
		if !b {
			return nil, false
		}
		providerBlock := module.NewProviderBlock()
		providerBlock.Provider = requiredProvider.Name
		providerBlock.Name = requiredProvider.Name
		providerBlock.Cache = "1d"
		providerBlock.MaxGoroutines = pointer.ToUInt64Pointer(50)
		providerBlock.ProvidersConfigYamlString = configuration
		providersBlock = append(providersBlock, providerBlock)
	}
	return providersBlock, true
}

// run provider & get it's init configuration
func (x *InitCommandExecutor) getProviderInitConfiguration(ctx context.Context, localProviderManager *local_providers_manager.LocalProvidersManager, plan *planner.ProviderInstallPlan) (string, bool) {

	// start & get information
	cli_ui.Infof("begin init provider %s", plan.String())

	// Find the local path of the provider
	localProvider := &local_providers_manager.LocalProvider{
		Provider: plan.Provider,
	}
	installed, d := localProviderManager.IsProviderInstalled(ctx, localProvider)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return "", false
	}
	if !installed {
		cli_ui.Errorf("provider %s not installed, can not exec fetch for it", plan.String())
		return "", false
	}

	// Find the local installation location of the provider
	localProviderMeta, d := localProviderManager.Get(ctx, localProvider)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return "", false
	}

	// Start provider
	plug, err := plugin.NewManagedPlugin(localProviderMeta.ExecutableFilePath, plan.Name, plan.Version, "", nil)
	if err != nil {
		cli_ui.Errorf("start provider %s at %s failed: %s", plan.String(), localProvider.ExecutableFilePath, err.Error())
		return "", false
	}
	// Close the provider at the end of the method execution
	defer plug.Close()

	cli_ui.Errorf("start provider %s success", plan.String())

	// Database connection option
	storageOpt := postgresql_storage.NewPostgresqlStorageOptions(x.options.DSN)
	providerBlock := module.NewProviderBlock()
	providerBlock.Name = plan.Name
	dbSchema := pgstorage.GetSchemaKey(plan.Name, plan.Version, providerBlock)
	pgstorage.WithSearchPath(dbSchema)(storageOpt)
	opt, err := json.Marshal(storageOpt)
	if err != nil {
		cli_ui.Errorf("json marshal postgresql options error: %s", err.Error())
		return "", false
	}

	// Get the lock first
	storage, d := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, storageOpt)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return "", false
	}
	lockId := "selefra-fetch-lock"
	ownerId := utils.BuildOwnerId()
	tryTimes := 0
	for {

		cli_ui.Infof("provider %s, schema %s, owner = %s, try get fetch lock...", plan.String(), dbSchema, ownerId)

		tryTimes++
		err := storage.Lock(ctx, lockId, ownerId)
		if err != nil {
			cli_ui.Errorf("provider %s, schema %s, owner = %s, get fetch lock error: %s, will sleep & retry, tryTimes = %d\n", plan.String(), dbSchema, ownerId, err.Error(), tryTimes)
		} else {
			cli_ui.Infof("provider %s, schema %s, owner = %s, get fetch lock success\n", plan.String(), dbSchema, ownerId)
			break
		}
		time.Sleep(time.Second * 10)
	}
	defer func() {
		for tryTimes := 0; tryTimes < 10; tryTimes++ {
			err := storage.UnLock(ctx, lockId, ownerId)
			if err != nil {
				cli_ui.Errorf("provider %s, schema %s, owner = %s, release fetch lock error: %s, will sleep & retry, tryTimes = %d", plan.String(), dbSchema, ownerId, err.Error(), tryTimes)
			} else {
				cli_ui.Infof("provider %s, schema %s, owner = %s, release fetch lock success", plan.String(), dbSchema, ownerId)
				break
			}
		}
	}()

	// Initialize the provider
	pluginProvider := plug.Provider()
	var providerYamlConfiguration string = module.GetDefaultProviderConfigYamlConfiguration(plan.Name, plan.Version)

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
		cli_ui.Errorf("start provider failed: %s", err.Error())
		return "", false
	}
	if err := cli_ui.PrintDiagnostics(providerInitResponse.Diagnostics); err != nil {
		return "", false
	}

	cli_ui.Infof("provider %s init success", plan.String())

	// Get information about the started provider
	information, err := pluginProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
	if err != nil {
		cli_ui.Errorf("provider %s, schema %s, get provider information failed: %s", plan.String(), dbSchema, err.Error())
		return "", false
	}

	return information.DefaultConfigTemplate, true
}

// ------------------------------------------------- --------------------------------------------------------------------
