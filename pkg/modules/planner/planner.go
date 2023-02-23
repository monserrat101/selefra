package planner

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
)

// Planner 表示一个可以生成计划的规划器
type Planner[T any] interface {

	// Name 规划器的名字
	Name() string

	// MakePlan 制定一个计划
	MakePlan(ctx context.Context) (T, *schema.Diagnostics)
}
