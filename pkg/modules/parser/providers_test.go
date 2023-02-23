package parser

import (
	"fmt"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlFileToModuleParser_parseProvidersBlock(t *testing.T) {
	module, diagnostics := NewYamlFileToModuleParser("./test_data/test_parse_providers/modules.yaml").Parse()
	assert.False(t, utils.HasError(diagnostics))
	assert.NotNil(t, module.ProvidersBlock)

	providerBlock := module.ProvidersBlock[0]
	assert.NotEmpty(t, providerBlock.GetNodeLocation("name.key").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("name.value").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("name").ReadSourceString())

	assert.NotEmpty(t, providerBlock.GetNodeLocation("cache.key").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("cache.value").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("cache").ReadSourceString())

	assert.NotEmpty(t, providerBlock.GetNodeLocation("resources.key").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("resources.value").ReadSourceString())
	assert.NotEmpty(t, providerBlock.GetNodeLocation("resources").ReadSourceString())

	for i := 0; i < len(providerBlock.Resources); i++ {
		assert.NotEmpty(t, providerBlock.GetNodeLocation(fmt.Sprintf("resources[%d]", i)).ReadSourceString())
		assert.NotEmpty(t, providerBlock.GetNodeLocation(fmt.Sprintf("resources[%d].value", i)).ReadSourceString())
	}
}
