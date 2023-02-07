package query

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/client"
	"testing"
)

func createCtxAndClient(cof config.RootConfig, required *config.ProviderRequired, cp config.ProviderConfig) (context.Context, *client.Client, error) {
	uid, _ := uuid.NewUUID()
	ctx := context.Background()
	c, e := client.CreateClientFromConfig(ctx, &cof.Selefra, uid, required, cp)
	if e != nil {
		ui.Errorln(e)
		return nil, nil, e
	}
	return ctx, c, nil
}

func TestCreateColumnsSuggest(t *testing.T) {
	ctx := context.Background()
	global.Init("go_test", global.WithWorkspace("../../tests/workspace/offline"))
	cof, err := config.GetConfig()
	if err != nil {
		ui.Errorln(err)
	}
	for i := range cof.Selefra.Providers {
		confs, err := tools.GetProviders(cof, cof.Selefra.Providers[i].Name)
		if err != nil {
			ui.Errorln(err.Error())
		}
		for _, conf := range confs {
			var cp config.ProviderConfig
			err := json.Unmarshal([]byte(conf), &cp)
			if err != nil {
				ui.Errorln(err.Error())
				continue
			}
			//ctx, c, err := createCtxAndClient(*cof, cof.Selefra.Providers[i], cp)
			//if err != nil {
			//	t.Error(err)
			//}
			sto, _ := pgstorage.Storage(ctx)
			columns := CreateColumnsSuggest(ctx, sto)
			if columns == nil {
				t.Error("Columns is nil")
			}
		}
	}
}

func TestCreateTablesSuggest(t *testing.T) {
	ctx := context.Background()
	global.Init("go_test", global.WithWorkspace("../../tests/workspace/offline"))
	cof, err := config.GetConfig()
	if err != nil {
		ui.Errorln(err)
	}
	for i := range cof.Selefra.Providers {
		confs, err := tools.GetProviders(cof, cof.Selefra.Providers[i].Name)
		if err != nil {
			ui.Errorln(err.Error())
		}
		for _, conf := range confs {
			var cp config.ProviderConfig
			err := json.Unmarshal([]byte(conf), &cp)
			if err != nil {
				ui.Errorln(err.Error())
				continue
			}
			sto, _ := pgstorage.Storage(ctx)
			tables := CreateTablesSuggest(ctx, sto)
			if tables == nil {
				t.Error("Tables is nil")
			}
		}
	}
}
