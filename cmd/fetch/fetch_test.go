package fetch

import (
	"testing"
)

//import (
//	"context"
//	"github.com/selefra/selefra/cmd/tools"
//	"github.com/selefra/selefra/config"
//	"github.com/selefra/selefra/global"
//	"testing"
//)
//
//func TestFetch(t *testing.T) {
//	ctx := context.Background()
//	global.Init("fetch", global.WithWorkspace("../../tests/workspace/offline"))
//	bootstrap, err := config.GetConfig()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for _, p := range bootstrap.Selefra.ProviderDecls {
//		confs, err := tools.ProviderConfigStrs(bootstrap, p.Name)
//		if err != nil {
//			t.Fatal(err)
//		}
//		for _, conf := range confs {
//			err = Fetch(ctx, p, conf)
//			if err != nil {
//				t.Error(err)
//			}
//		}
//	}
//}

func TestFetch(t *testing.T) {
	projectWorkspace := "./test_data/test_fetch_module"
	downloadWorkspace := "./test_download"
	Fetch(projectWorkspace, downloadWorkspace)
}
