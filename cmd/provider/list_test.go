package provider

import (
	"github.com/selefra/selefra/global"
	"testing"
)

func TestList(t *testing.T) {
	global.Init("TestList", global.WithWorkspace("../../tests/workspace/offline"))
	err := list()
	if err != nil {
		t.Error(err)
	}
}
