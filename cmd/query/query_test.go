package query

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestNewQueryClient(t *testing.T) {
	ctx := context.Background()
	global.Init("query", global.WithWorkspace("../../tests/workspace/offline"))

	queryClient, _ := NewQueryClient(ctx)
	if queryClient == nil {
		t.Error("queryClient is nil")
	}
}

func TestNewQueryClientOnline(t *testing.T) {
	ctx := context.Background()
	global.Init("query", global.WithWorkspace("../../tests/workspace/online"))
	global.SetToken("4fe8ed36488c479d0ba7292fe09a4132")
	global.SERVER = "dev-api.selefra.io"

	queryClient, _ := NewQueryClient(ctx)
	if queryClient == nil {
		t.Error("queryClient is nil")
	}
}
