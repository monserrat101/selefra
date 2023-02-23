package provider

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

//import (
//	"context"
//	"github.com/selefra/selefra/global"
//	"testing"
//)
//
//func TestRemove(t *testing.T) {
//	global.Init("TestRemove", global.WithWorkspace("../../tests/workspace/offline"))
//	err := Remove([]string{"aws"})
//	if err != nil {
//		t.Error(err)
//	}
//	err = install(context.Background(), []string{"aws@latest"})
//	if err != nil {
//		t.Error(err)
//	}
//}

func TestRemove(t *testing.T) {
	err := Remove(context.Background(), "./test_download", "aws@v0.0.1")
	assert.Nil(t, err)
}
