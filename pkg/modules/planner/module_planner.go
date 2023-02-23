package planner

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/modules/module"
)

// ------------------------------------------------- --------------------------------------------------------------------

// MakeModuleQueryPlan Generate an execution plan for the module
func MakeModuleQueryPlan(ctx context.Context, module *module.Module, tableToProviderMap map[string]string) (*ModulePlan, *schema.Diagnostics) {
	return NewModulePlanner(module, tableToProviderMap).MakePlan(ctx)
}

// ------------------------------------------------- --------------------------------------------------------------------

// ModulePlan Represents the execution plan of a module
type ModulePlan struct {

	// Which module is this execution plan generated for
	*module.Module

	// 模块级别的作用域
	ModuleScope *Scope

	// 子模块的执行计划
	SubModulesPlan []*ModulePlan

	// 这个模块下的rule的执行规划
	RulesPlan []*RulePlan
}

//// ------------------------------------------------- --------------------------------------------------------------------
//
//// RootModulePlan The execution plan of the root module
//type RootModulePlan struct {
//
//	// The root module's execution plan is also a module execution plan
//	*ModulePlan
//
//	// The provider pull plan for all the following modules is extracted to the root module level
//	ProviderFetchPlanSlice []*ProviderFetchPlan
//}
//

// ------------------------------------------------- --------------------------------------------------------------------

// ModulePlanner Used to generate an execution plan for a module
type ModulePlanner struct {
	module             *module.Module
	tableToProviderMap map[string]string
}

var _ Planner[*ModulePlan] = &ModulePlanner{}

func NewModulePlanner(module *module.Module, tableToProviderMap map[string]string) *ModulePlanner {
	return &ModulePlanner{
		module:             module,
		tableToProviderMap: tableToProviderMap,
	}
}

func (x *ModulePlanner) Name() string {
	return "module-planner"
}

func (x *ModulePlanner) MakePlan(ctx context.Context) (*ModulePlan, *schema.Diagnostics) {
	return x.buildModulePlanner(ctx, x.module, NewScope(), x.tableToProviderMap)
}

func (x *ModulePlanner) buildModulePlanner(ctx context.Context, module *module.Module, scope *Scope, tableToProviderMap map[string]string) (*ModulePlan, *schema.Diagnostics) {
	diagnostics := schema.NewDiagnostics()
	modulePlan := &ModulePlan{
		Module: module,
	}

	// Generate an execution plan for the rules in the module
	for _, ruleBlock := range module.RulesBlock {
		rulePlan, d := NewRulePlanner(ruleBlock, scope, tableToProviderMap).MakePlan(ctx)
		rulePlan.ModulePlan = modulePlan
		rulePlan.Module = module
		if diagnostics.Add(d).HasError() {
			return nil, diagnostics
		}
		modulePlan.RulesPlan = append(modulePlan.RulesPlan, rulePlan)
	}

	// 为子模块生成执行计划
	for _, subModule := range module.SubModules {

		// 子模块的作用域
		subModuleScope := NewScope()
		// 子模块的作用域继承当前模块
		subModuleScope.Extend(scope)
		// 同时子模块可能有一些初始化的变量
		// TODO

		subModulePlan, d := x.buildModulePlanner(ctx, subModule, subModuleScope, tableToProviderMap)
		if diagnostics.AddDiagnostics(d).HasError() {
			return nil, diagnostics
		}
		modulePlan.SubModulesPlan = append(modulePlan.SubModulesPlan, subModulePlan)
	}

	return modulePlan, diagnostics
}

// ------------------------------------------------- --------------------------------------------------------------------
