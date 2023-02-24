package planner

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/storage/pgstorage"
)

// ------------------------------------------------- --------------------------------------------------------------------

type ProvidersFetchPlan []*ProviderFetchPlan

func (x ProvidersFetchPlan) BuildProviderContextMap(ctx context.Context, DSN string) (map[string][]*ProviderContext, *schema.Diagnostics) {
	diagnostics := schema.NewDiagnostics()
	m := make(map[string][]*ProviderContext, 0)
	for _, plan := range x {

		databaseSchema := pgstorage.GetSchemaKey(plan.Name, plan.Version, plan.ProviderConfigurationBlock)
		options := postgresql_storage.NewPostgresqlStorageOptions(DSN)
		options.SearchPath = databaseSchema

		databaseStorage, d := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, options)
		if diagnostics.AddDiagnostics(d).HasError() {
			return nil, diagnostics
		}

		providerContext := &ProviderContext{
			ProviderName:          plan.Name,
			ProviderVersion:       plan.Version,
			Schema:                databaseSchema,
			Storage:               databaseStorage,
			ProviderConfiguration: plan.ProviderConfigurationBlock,
		}
		m[plan.Name] = append(m[plan.Name], providerContext)
	}

	return m, diagnostics
}

// ProviderContext Ready execution strategy
type ProviderContext struct {

	// Which provider is it?
	ProviderName string

	// Which version
	ProviderVersion string

	// The database stored to
	Schema string

	// A connection to a database instance
	Storage storage.Storage

	// The provider configuration block
	ProviderConfiguration *module.ProviderBlock
}

// ------------------------------------------------- --------------------------------------------------------------------

const (
	DefaultMaxGoroutines = uint64(100)
)

// ProviderFetchPlan Indicates the pull plan of a provider
type ProviderFetchPlan struct {
	*ProviderInstallPlan

	// provider Configuration information used for fetching
	ProviderConfigurationBlock *module.ProviderBlock
}

func NewProviderFetchPlan(providerName, providerVersion string, providerBlock *module.ProviderBlock) *ProviderFetchPlan {
	return &ProviderFetchPlan{
		ProviderInstallPlan: &ProviderInstallPlan{
			Provider: registry.NewProvider(providerName, providerVersion),
		},
		ProviderConfigurationBlock: providerBlock,
	}
}

// GetProvidersConfigYamlString Obtain the configuration file for running the Provider
func (x *ProviderFetchPlan) GetProvidersConfigYamlString() string {
	if x.ProviderConfigurationBlock != nil {
		return x.ProviderConfigurationBlock.ProvidersConfigYamlString
	}
	return ""
}

// GetNeedPullTablesName Gets which tables to pull when pulling
func (x *ProviderFetchPlan) GetNeedPullTablesName() []string {
	tables := make([]string, 0)
	if x.ProviderConfigurationBlock != nil {
		tables = x.ProviderConfigurationBlock.Resources
	}
	if len(tables) == 0 {
		tables = append(tables, provider.AllTableNameWildcard)
	}
	return tables
}

// GetMaxGoroutines How many concurrency is used to pull the table data
func (x *ProviderFetchPlan) GetMaxGoroutines() uint64 {
	if x.ProviderConfigurationBlock != nil && x.ProviderConfigurationBlock.MaxGoroutines != nil {
		return *x.ProviderConfigurationBlock.MaxGoroutines
	} else {
		return DefaultMaxGoroutines
	}
}

// ------------------------------------------------- --------------------------------------------------------------------

// ProviderFetchPlannerOptions This parameter is required when creating the provider execution plan
type ProviderFetchPlannerOptions struct {

	// Which module is the execution plan being generated for
	Module *module.Module

	// Provider version that wins the vote
	ProviderVersionVoteWinnerMap map[string]string
}

// ------------------------------------------------- --------------------------------------------------------------------

type ProviderFetchPlanner struct {
	options *ProviderFetchPlannerOptions
}

var _ Planner[ProvidersFetchPlan] = &ProviderFetchPlanner{}

func NewProviderFetchPlanner(options *ProviderFetchPlannerOptions) *ProviderFetchPlanner {
	return &ProviderFetchPlanner{
		options: options,
	}
}

func (x *ProviderFetchPlanner) Name() string {
	return "provider-fetch-planner"
}

func (x *ProviderFetchPlanner) MakePlan(ctx context.Context) (ProvidersFetchPlan, *schema.Diagnostics) {
	return x.expandByConfiguration()
}

// Expand to multiple tasks based on the configuration
func (x *ProviderFetchPlanner) expandByConfiguration() ([]*ProviderFetchPlan, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()
	providerFetchPlanSlice := make([]*ProviderFetchPlan, 0)

	if x.options.Module.SelefraBlock == nil {
		return nil, diagnostics.AddErrorMsg("module %s must have selefra block for make fetch plan", x.options.Module.BuildFullName())
	} else if len(x.options.Module.SelefraBlock.RequireProvidersBlock) == 0 {
		return nil, diagnostics.AddErrorMsg("module %s selefra block not have providers block", x.options.Module.BuildFullName())
	}

	// Start a task for those that have a task written, some join by fetch start rule
	providerNamePlanCountMap := make(map[string]int, 0)
	nameToProviderMap := x.options.Module.SelefraBlock.RequireProvidersBlock.BuildNameToProviderBlockMap()
	for _, providerBlock := range x.options.Module.ProvidersBlock {
		requiredProviderBlock, exists := nameToProviderMap[providerBlock.Provider]
		if !exists {
			// selefra.providers block not found that name in providers[index] configuration
			errorTips := fmt.Sprintf("provider name %s not found", providerBlock.Provider)
			diagnostics.AddErrorMsg(module.RenderErrorTemplate(errorTips, providerBlock.GetNodeLocation("")))
		} else if providerWinnerVersion, exists := x.options.ProviderVersionVoteWinnerMap[requiredProviderBlock.Source]; exists {
			// Start a plan for the provider
			providerNamePlanCountMap[requiredProviderBlock.Source]++
			providerFetchPlanSlice = append(providerFetchPlanSlice, NewProviderFetchPlan(requiredProviderBlock.Source, providerWinnerVersion, providerBlock))
		} else {
			errorTips := fmt.Sprintf("provider version %s not found", requiredProviderBlock.Source)
			diagnostics.AddErrorMsg(module.RenderErrorTemplate(errorTips, requiredProviderBlock.GetNodeLocation("version")))
		}
	}
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	// See if there is another project that has not been activated, and if there is, start a pull plan for it as well
	for providerName, providerVersion := range x.options.ProviderVersionVoteWinnerMap {
		if providerNamePlanCountMap[providerName] > 0 {
			continue
		}
		providerFetchPlanSlice = append(providerFetchPlanSlice, NewProviderFetchPlan(providerName, providerVersion, nil))
	}

	return providerFetchPlanSlice, diagnostics
}

// ------------------------------------------------- --------------------------------------------------------------------
