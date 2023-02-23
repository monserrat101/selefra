package planner

import (
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/utils"
)

// Scope Used to represent the moduleScope of a module
type Scope struct {

	// Variable in moduleScope
	variablesMap map[string]any

	// The configuration information in moduleScope
	providerConfigBlockSlice []*module.ProviderBlock
}

func ExtendScope(scope *Scope) *Scope {
	subScope := NewScope()
	subScope.Extend(scope)
	return subScope
}

func NewScope() *Scope {
	return &Scope{
		variablesMap: make(map[string]any),
	}
}

// Extend moduleScope inheritance
func (x *Scope) Extend(scope *Scope) {
	for key, value := range scope.variablesMap {
		if _, exists := x.variablesMap[key]; exists {
			continue
		}
		x.variablesMap[key] = value
	}
}

// Clone Make a copy of the current moduleScope
func (x *Scope) Clone() *Scope {

	newVariablesMap := make(map[string]any)
	for key, value := range x.variablesMap {
		newVariablesMap[key] = value
	}

	return &Scope{
		variablesMap:             newVariablesMap,
		providerConfigBlockSlice: x.providerConfigBlockSlice,
	}
}

// GetVariable Gets the value of a variable
func (x *Scope) GetVariable(variableName string) (any, bool) {
	value, exists := x.variablesMap[variableName]
	return value, exists
}

// DeclareVariable Declare a variable
func (x *Scope) DeclareVariable(variableName string, variableValue any) any {
	oldValue := x.variablesMap[variableName]
	x.variablesMap[variableName] = variableValue
	return oldValue
}

// DeclareVariableIfNotExists Declared only if the variable does not exist
func (x *Scope) DeclareVariableIfNotExists(variableName string, variableValue any) bool {
	if _, exists := x.variablesMap[variableName]; exists {
		return false
	}
	x.variablesMap[variableName] = variableValue
	return true
}

// DeclareVariables Batch declaration variable
func (x *Scope) DeclareVariables(variablesMap map[string]any) {
	for variableName, variableValue := range variablesMap {
		x.variablesMap[variableName] = variableValue
	}
}

// RenderingTemplate Rendering the template using the moduleScope of the current module
func (x *Scope) RenderingTemplate(templateName, templateString string) (string, error) {
	return utils.RenderingTemplate(templateName, templateString, x.variablesMap)
}
