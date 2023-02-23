package parser

import (
	"fmt"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlFileToModuleParser_Parse(t *testing.T) {
	module, diagnostics := NewYamlFileToModuleParser("./test_data/test.yaml").Parse()
	assert.False(t, utils.HasError(diagnostics))
	t.Log(diagnostics.ToString())
	t.Log(module)

	location := module.RulesBlock[0].MetadataBlock.GetNodeLocation("tags[0].value")
	s := location.ReadSourceString()
	fmt.Println(s)
}
