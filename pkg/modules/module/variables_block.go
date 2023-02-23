package module

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/reflect_util"
)

// ------------------------------------------------- --------------------------------------------------------------------

// VariablesBlock One of the root-level code blocks
type VariablesBlock []*VariableBlock

var _ Block = (*VariablesBlock)(nil)
var _ MergableBlock[VariablesBlock] = (*VariablesBlock)(nil)

func (x VariablesBlock) Merge(other VariablesBlock) (VariablesBlock, *schema.Diagnostics) {

	diagnostics := schema.NewDiagnostics()

	variableKeySet := make(map[string]struct{}, 0)
	mergedVariables := make(VariablesBlock, 0)

	// merge self
	for _, variableBlock := range x {
		if _, exists := variableKeySet[variableBlock.Key]; exists {
			// TODO
			diagnostics.AddErrorMsg("variables block merge failed, find same name variable %s", variableBlock.Key)
			continue
		}
		mergedVariables = append(mergedVariables, variableBlock)
		variableKeySet[variableBlock.Key] = struct{}{}
	}

	// merge other
	for _, variableBlock := range other {
		if _, exists := variableKeySet[variableBlock.Key]; exists {
			// TODO
			diagnostics.AddErrorMsg("variables block merge failed, find same name variable %s", variableBlock.Key)
			continue
		}
		mergedVariables = append(mergedVariables, variableBlock)
		variableKeySet[variableBlock.Key] = struct{}{}
	}

	return mergedVariables, diagnostics
}

func (x VariablesBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	variableKeySet := make(map[string]struct{}, 0)
	for _, variableBlock := range x {
		if _, exists := variableKeySet[variableBlock.Key]; exists {
			// TODO
			diagnostics.AddErrorMsg("variables block check failed, find same name variable")
			continue
		}
		variableKeySet[variableBlock.Key] = struct{}{}
		diagnostics.AddDiagnostics(variableBlock.Check(module, validatorContext))
	}

	return diagnostics
}

func (x VariablesBlock) IsEmpty() bool {
	return len(x) == 0
}

func (x VariablesBlock) GetNodeLocation(selector string) *NodeLocation {
	panic("not supported")
}

func (x VariablesBlock) SetNodeLocation(selector string, nodeLocation *NodeLocation) error {
	panic("not supported")
}

// ------------------------------------------------- --------------------------------------------------------------------

// VariableBlock Used to declare a variable
type VariableBlock struct {

	// Name of a variable
	Key string `yaml:"key" json:"key"`

	// The default value of the variable
	Default any `yaml:"default" json:"default"`

	// A description of this variable
	Description string `yaml:"description" json:"description"`

	// Who is the author of the variable? What the hell is this?
	Author string `yaml:"author" json:"author"`

	*LocatableImpl `yaml:"-"`
}

var _ Block = &VariableBlock{}

func NewVariableBlock() *VariableBlock {
	return &VariableBlock{
		LocatableImpl: NewLocatableImpl(),
	}
}

func (x *VariableBlock) Check(module *Module, validatorContext *ValidatorContext) *schema.Diagnostics {

	diagnostics := schema.NewDiagnostics()

	if x.Key == "" {
		// TODO build error message
		diagnostics.AddErrorMsg("variable block .key must be specified")
	}

	if reflect_util.IsNil(x.Default) {
		// TODO build error message
		diagnostics.AddErrorMsg("variable block .default must be specified")
	}

	return diagnostics
}

func (x *VariableBlock) IsEmpty() bool {
	return x.Key == "" &&
		reflect_util.IsNil(x.Default) &&
		x.Description == "" &&
		x.Author == ""
}

// ------------------------------------------------- --------------------------------------------------------------------
