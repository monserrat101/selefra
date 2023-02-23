package planner

import (
	"context"
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

type ProviderFetchPlanner struct {
	module                       *module.Module
	providerVersionVoteWinnerMap map[string]string
}

var _ Planner[ProvidersFetchPlan] = &ProviderFetchPlanner{}

func NewProviderFetchPlanner(module *module.Module, providerVersionVoteWinnerMap map[string]string) *ProviderFetchPlanner {
	return &ProviderFetchPlanner{
		module:                       module,
		providerVersionVoteWinnerMap: providerVersionVoteWinnerMap,
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

	// Start a task for those that have a task written
	providerNamePlanCountMap := make(map[string]int, 0)
	providerNameMap := x.module.SelefraBlock.RequireProvidersBlock.ToNameMap()
	for _, providerBlock := range x.module.ProvidersBlock {
		block, exists := providerNameMap[providerBlock.Provider]
		if !exists {
			// TODO provider name not found
			diagnostics.AddErrorMsg("provider name %s not found", providerBlock.Provider)
		} else if providerWinnerVersion, exists := x.providerVersionVoteWinnerMap[block.Source]; exists {
			// Start a plan for the provider
			providerNamePlanCountMap[block.Source]++
			providerFetchPlanSlice = append(providerFetchPlanSlice, NewProviderFetchPlan(block.Source, providerWinnerVersion, providerBlock))
		} else {
			// TODO  provider version not found
			diagnostics.AddErrorMsg("provider version %s not found", block.Source)
		}
	}
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	// See if there is another project that has not been activated, and if there is, start a pull plan for it as well
	for providerName, providerVersion := range x.providerVersionVoteWinnerMap {
		if providerNamePlanCountMap[providerName] > 0 {
			continue
		}
		providerFetchPlanSlice = append(providerFetchPlanSlice, NewProviderFetchPlan(providerName, providerVersion, nil))
	}

	return providerFetchPlanSlice, diagnostics
}

// ------------------------------------------------- --------------------------------------------------------------------
