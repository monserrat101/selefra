package provider

import (
	"context"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestRemoveOnline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	global.Init("TestRemoveOnline", global.WithWorkspace("../../tests/workspace/online"))
	global.SetToken("4fe8ed36488c479d0ba7292fe09a4132")
	global.SERVER = "dev-api.selefra.io"
	err := Remove([]string{"aws"})
	if err != nil {
		t.Error(err)
	}
	err = install(context.Background(), []string{"aws@latest"})
	if err != nil {
		t.Error(err)
	}
}
