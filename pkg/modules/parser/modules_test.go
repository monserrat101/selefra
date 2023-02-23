package parser

import (
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlFileToModuleParser_parseModulesBlock(t *testing.T) {
	module, diagnostics := NewYamlFileToModuleParser("./test_data/test_parse_modules/modules.yaml").Parse()
	assert.False(t, utils.HasError(diagnostics))
	assert.NotNil(t, module.ModulesBlock)

	moduleBlock := module.ModulesBlock[1]
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("name.key").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("name.value").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("name").ReadSourceString())

	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input.key").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input.value").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input").ReadSourceString())

	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input.name.key").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input.name.value").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("input.name").ReadSourceString())

	assert.NotEmpty(t, moduleBlock.GetNodeLocation("uses.key").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("uses.value").ReadSourceString())
	assert.NotEmpty(t, moduleBlock.GetNodeLocation("uses").ReadSourceString())

}
