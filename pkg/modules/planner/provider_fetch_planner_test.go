package planner

import (
	"context"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderFetchPlanner_MakePlan(t *testing.T) {

	rootModule := module.NewModule()
	rootModule.SelefraBlock = module.NewSelefraBlock()
	rootModule.SelefraBlock.RequireProvidersBlock = []*module.RequireProviderBlock{
		{
			Name:    "aws",
			Source:  "aws",
			Version: "latest",
		},
		{
			Name:    "gcp",
			Source:  "gcp",
			Version: "latest",
		},
	}
	rootModule.ProvidersBlock = []*module.ProviderBlock{
		{
			Name:          "aws-001",
			Provider:      "aws",
			MaxGoroutines: pointer.ToUInt64Pointer(10),
		},
		{
			Name:          "aws-002",
			Provider:      "aws",
			MaxGoroutines: pointer.ToUInt64Pointer(30),
		},
	}
	versionWinnerMap := map[string]string{
		"aws": "v0.0.1",
		"gcp": "v0.0.1",
	}
	plan, diagnostics := NewProviderFetchPlanner(rootModule, versionWinnerMap).MakePlan(context.Background())
	assert.False(t, utils.HasError(diagnostics))
	assert.Equal(t, 3, len(plan))
}
