package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestListOnline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	global.Init("TestListOnline", global.WithWorkspace("../../tests/workspace/online"))
	global.SetToken("4fe8ed36488c479d0ba7292fe09a4132")
	global.SERVER = "dev-api.selefra.io"
	err := list()
	if err != nil {
		t.Error(err)
	}
}
