package test

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestTestFuncOnline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	global.Init("TestTestFuncOnline", global.WithWorkspace("../../tests/workspace/online"))
	global.SetToken("4fe8ed36488c479d0ba7292fe09a4132")
	global.SERVER = "dev-api.selefra.io"

	ctx := context.Background()
	err := TestFunc(ctx)
	if err != nil {
		t.Error(err)
	}
}
