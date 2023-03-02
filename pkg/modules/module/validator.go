package module

// ------------------------------------------------- --------------------------------------------------------------------

// ValidatorContext Some global context information stored during validation
type ValidatorContext struct {

	// Global collection of rule ids
	RulesIdSet map[string]*RuleBlock

	// All module names, if there are module names such as the same name should be able to check out
	ModuleNameSet map[string]*ModuleBlock
}

// NewValidatorContext Create a validation context
func NewValidatorContext() *ValidatorContext {
	return &ValidatorContext{
		RulesIdSet:    make(map[string]*RuleBlock),
		ModuleNameSet: make(map[string]*ModuleBlock),
	}
}

// AddRuleBlock Add rules to the validation context
func (x *ValidatorContext) AddRuleBlock(ruleBlock *RuleBlock) {
	if ruleBlock.MetadataBlock != nil {
		x.RulesIdSet[ruleBlock.MetadataBlock.Id] = ruleBlock
	}
}

// GetRuleBlockById Determine whether the given rule is in context
func (x *ValidatorContext) GetRuleBlockById(ruleId string) (*RuleBlock, bool) {
	ruleBlock, exists := x.RulesIdSet[ruleId]
	return ruleBlock, exists
}

// AddModuleBlock Adds the module to the current validator context
func (x *ValidatorContext) AddModuleBlock(moduleBlock *ModuleBlock) {
	x.ModuleNameSet[moduleBlock.Name] = moduleBlock
}

// GetModuleByName Gets the module in the validation context
func (x *ValidatorContext) GetModuleByName(moduleName string) (*ModuleBlock, bool) {
	moduleBlock, exists := x.ModuleNameSet[moduleName]
	return moduleBlock, exists
}

// ------------------------------------------------- --------------------------------------------------------------------

const CheckIdentityErrorMsg = "only allow \"a-z,A-Z,0-9,_\" and can't start with a number"

func checkIdentity(s string) bool {

	if len(s) == 0 {
		return false
	}

	// And you can't start with a number
	if s[0] >= '0' && s[0] <= '9' {
		return false
	}

	// Only the given character can be used
	for _, c := range s {
		isOk := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || (c == '_')
		if !isOk {
			return false
		}
	}
	return true
}

// ------------------------------------------------- --------------------------------------------------------------------
