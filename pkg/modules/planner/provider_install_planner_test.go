package planner

import (
	"context"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProviderInstallPlanner(t *testing.T) {

	// case 1: 能够选出一个明确的版本
	rootModule := randomModule("v0.0.1, v0.0.2")
	rootModule.SubModules = append(rootModule.SubModules, randomModule("v0.0.2, v0.0.3"))
	plan, diagnostics := NewProviderInstallPlanner(rootModule).MakePlan(context.Background())
	assert.False(t, utils.HasError(diagnostics))
	assert.Len(t, plan, 1)
	assert.Equal(t, "v0.0.2", plan[0].Version)

	//// case 2: 能够选出多个明确的版本
	rootModule = randomModule("v0.0.1, v0.0.2, v0.0.3")
	rootModule.SubModules = append(rootModule.SubModules, randomModule("v0.0.2, v0.0.3, v0.0.4"))
	plan, diagnostics = NewProviderInstallPlanner(rootModule).MakePlan(context.Background())
	assert.False(t, utils.HasError(diagnostics))
	assert.Len(t, plan, 1)
	assert.Equal(t, "v0.0.3", plan[0].Version)
	//
	//// case 3: 无法选出明确的版本
	rootModule = randomModule("v0.0.1, v0.0.2")
	rootModule.SubModules = append(rootModule.SubModules, randomModule("v0.0.3, v0.0.4"))
	plan, diagnostics = NewProviderInstallPlanner(rootModule).MakePlan(context.Background())
	assert.True(t, utils.HasError(diagnostics))
	t.Log(diagnostics.ToString())
	assert.Len(t, plan, 0)

}
