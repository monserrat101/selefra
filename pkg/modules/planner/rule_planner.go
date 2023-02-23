package planner

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/modules/module"
)

// ------------------------------------------------- --------------------------------------------------------------------

type RulePlan struct {

	// The execution plan of the module to which it is associated
	ModulePlan *ModulePlan

	// The module to which it is associated
	Module *module.Module

	// Is the execution plan for which block
	*module.RuleBlock

	// Which provider is the rule bound to? Currently, a rule can be bound to only one provider
	BindingProviderName string

	// Render a good rule - bound Query
	Query string

	// Which tables are used in this Query
	BindingTables []string

	RuleScope *Scope
}

func (x *RulePlan) String() string {
	if x.MetadataBlock != nil {
		return x.Name + ":" + x.MetadataBlock.Id
	} else {
		return x.Name
	}
}

// ------------------------------------------------- --------------------------------------------------------------------

// MakeRulePlan Plan the execution of the rule
func MakeRulePlan(ctx context.Context, ruleBlock *module.RuleBlock, moduleScope *Scope, tableToProviderMap map[string]string) (*RulePlan, *schema.Diagnostics) {
	return NewRulePlanner(ruleBlock, moduleScope, tableToProviderMap).MakePlan(ctx)
}

// ------------------------------------------------- --------------------------------------------------------------------

// RulePlanner An enforcement plan for this rule
type RulePlanner struct {
	ruleBlock          *module.RuleBlock
	moduleScope        *Scope
	tableToProviderMap map[string]string
}

var _ Planner[*RulePlan] = &RulePlanner{}

func (x *RulePlanner) Name() string {
	return "rule-planner"
}

func NewRulePlanner(ruleBlock *module.RuleBlock, moduleScope *Scope, tableToProviderMap map[string]string) *RulePlanner {
	return &RulePlanner{
		ruleBlock:          ruleBlock,
		moduleScope:        moduleScope,
		tableToProviderMap: tableToProviderMap,
	}
}

func (x *RulePlanner) MakePlan(ctx context.Context) (*RulePlan, *schema.Diagnostics) {
	diagnostics := schema.NewDiagnostics()
	ruleScope := ExtendScope(x.moduleScope)
	query, err := ruleScope.RenderingTemplate(x.ruleBlock.Query, x.ruleBlock.Query)
	if err != nil {
		location := x.ruleBlock.GetNodeLocation("query._value")
		report := module.RenderErrorTemplate(fmt.Sprintf("rendering template error: %s", err.Error()), location)
		return nil, diagnostics.AddErrorMsg(report)
	}
	bindingProviders, bindingTables := x.extractBinding(query, x.tableToProviderMap)
	if len(bindingProviders) != 1 {
		// TODO
		return nil, diagnostics.AddErrorMsg("rule must and only binding one provider: %s", x.ruleBlock.Query)
	}
	return &RulePlan{
		RuleBlock:           x.ruleBlock,
		BindingProviderName: bindingProviders[0],
		Query:               query,
		BindingTables:       bindingTables,
		RuleScope:           ruleScope,
	}, diagnostics
}

// Extract the names of the tables it uses from the rendered rule Query
func (x *RulePlanner) extractBinding(query string, tableToProviderMap map[string]string) (bindingProviders []string, bindingTables []string) {
	bindingProviderSet := make(map[string]struct{})
	bindingTableSet := make(map[string]struct{})
	inWord := false
	lastIndex := 0
	for index, c := range query {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= '0' && c <= '9' {
			if !inWord {
				inWord = true
				lastIndex = index
			}
		} else {
			if inWord {
				word := query[lastIndex:index]
				if providerName, exists := tableToProviderMap[word]; exists {
					bindingTableSet[word] = struct{}{}
					bindingProviderSet[providerName] = struct{}{}
				}
				inWord = false
			}
		}
	}

	for providerName := range bindingProviderSet {
		bindingProviders = append(bindingProviders, providerName)
	}
	for tableName := range bindingTableSet {
		bindingTables = append(bindingTables, tableName)
	}
	return
}

// ------------------------------------------------- --------------------------------------------------------------------

//// Extracting the provider name from the table name used by the policy is an implicit association
//func (x *RulePlanner) extractImplicitProvider(tablesName []string) ([]string, *schema.Diagnostics) {
//	diagnostics := schema.NewDiagnostics()
//	providerNameSet := make(map[string]struct{}, 0)
//	for _, tableName := range tablesName {
//		split := strings.SplitN(tableName, "_", 2)
//		if len(split) != 2 {
//			diagnostics.AddErrorMsg("can not found implicit provider name from table name %s", tableName)
//		} else {
//			providerNameSet[split[0]] = struct{}{}
//		}
//	}
//	providerNameSlice := make([]string, 0)
//	for providerName := range providerNameSet {
//		providerNameSlice = append(providerNameSlice, providerName)
//	}
//	return providerNameSlice, diagnostics
//}
//
//// Extract the names of the tables it uses from the rendered rule Query
//func (x *RulePlanner) extractTableNameSliceFromRuleQuery(s string, whitelistWordSet map[string]string) []string {
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

// ------------------------------------------------- --------------------------------------------------------------------
