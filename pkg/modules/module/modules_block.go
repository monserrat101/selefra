package module

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
)

// ------------------------------------------------- --------------------------------------------------------------------

type ModulesBlock []*ModuleBlock

var _ Block = (*ModulesBlock)(nil)
var _ MergableBlock[ModulesBlock] = (*ModulesBlock)(nil)

func (x ModulesBlock) Merge(other ModulesBlock) (ModulesBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	moduleNameSet := make(map[string]struct{})
	mergedModules := make(ModulesBlock, 0)

	// merge myself
	for _, moduleBlock := range x {
		if _, exists := moduleNameSet[moduleBlock.Name]; exists {
			// TODO error message
			diagnostics.AddErrorMsg("merge modules block error, find same name module %s", moduleBlock.Name)
			continue
		}
		mergedModules = append(mergedModules, moduleBlock)
		moduleNameSet[moduleBlock.Name] = struct{}{}
	}

	// merge other
	for _, moduleBlock := range other {
		if _, exists := moduleNameSet[moduleBlock.Name]; exists {
			// TODO error message
			diagnostics.AddErrorMsg("merge modules block error, find same name module %s", moduleBlock.Name)
			continue
		}
		mergedModules = append(mergedModules, moduleBlock)
		moduleNameSet[moduleBlock.Name] = struct{}{}
	}

	return mergedModules, diagnostics
}

func (x ModulesBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {
	diagnostics := schema.NewDiagnostics()
	for _, moduleBlock := range x {
		diagnostics.AddDiagnostics(moduleBlock.Check(module, validatorContext))
	}
	return diagnostics
}

func (x ModulesBlock) IsEmpty() bool {
	return len(x) == 0
}

func (x ModulesBlock) GetNodeLocation(selector string) *NodeLocation {
	panic("not supported")
}

func (x ModulesBlock) SetNodeLocation(selector string, nodeLocation *NodeLocation) error {
	panic("not supported")
}

// ------------------------------------------------- --------------------------------------------------------------------

// ModuleBlock Used to represent a common element in the modules array
type ModuleBlock struct {

	// Module name
	Name string `yaml:"name" json:"name"`

	// What other modules are referenced by this module
	Uses []string `yaml:"uses" json:"uses"`

	// The module supports specifying some variables
	Input map[string]any `yaml:"input" json:"input"`

	*LocatableImpl `yaml:"-"`
}

var _ Block = &ModuleBlock{}

func NewModuleBlock() *ModuleBlock {
	return &ModuleBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

func (x *ModuleBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	if x.Name == "" {
		// TODO error message
		diagnostics.AddErrorMsg("module name can not be empty")
	}

	if len(x.Uses) == 0 {
		// // TODO error message
		diagnostics.AddErrorMsg("module uses can not be empty")
	}

	if len(x.Input) != 0 {
		// TODO error message
		diagnostics.AddDiagnostics(x.checkInput())
	}

	return diagnostics
}

func (x *ModuleBlock) checkInput() *schema.Diagnostics {
	// TODO
	return nil
}

func (x *ModuleBlock) IsEmpty() bool {
	return x.Name == "" && len(x.Uses) == 0 && len(x.Input) == 0
}

// ------------------------------------------------- --------------------------------------------------------------------
