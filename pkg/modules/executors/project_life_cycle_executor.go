package executors

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
)

// ------------------------------------------------ ---------------------------------------------------------------------

type ProjectLifeApply struct {
}

// ------------------------------------------------ ---------------------------------------------------------------------

type ProjectLifeCycleExecutorOptions struct {
}

// ------------------------------------------------ ---------------------------------------------------------------------

const ProjectLifeCycleExecutorName = "project-life-cycle-executor"

// ProjectLifeCycleExecutor Used to fully run the entire project lifecycle
type ProjectLifeCycleExecutor struct {
	options *ProjectLifeCycleExecutorOptions
}

var _ Executor = &ProjectLifeCycleExecutor{}

func (x *ProjectLifeCycleExecutor) Name() string {
	return ProjectLifeCycleExecutorName
}

func (x *ProjectLifeCycleExecutor) Execute(ctx context.Context) *schema.Diagnostics {

	

}

// ------------------------------------------------ ---------------------------------------------------------------------
