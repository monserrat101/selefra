package fetch

import (
	"context"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"testing"
)

func TestFetch(t *testing.T) {
	ctx := context.Background()
	global.Init("fetch", global.WithWorkspace("../../tests/workspace/offline"))
	bootstrap, err := config.GetConfig()
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range bootstrap.Selefra.Providers {
		confs, err := tools.GetProviders(bootstrap, p.Name)
		if err != nil {
			t.Fatal(err)
		}
		for _, conf := range confs {
			err = Fetch(ctx, bootstrap, p, conf)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
