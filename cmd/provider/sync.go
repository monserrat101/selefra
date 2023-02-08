package provider

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/id_util"
	"github.com/selefra/selefra/cmd/fetch"
	"github.com/selefra/selefra/cmd/test"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/logger"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"path/filepath"
	"time"
)

type lockStruct struct {
	SchemaKey string
	Uuid      string
	Storage   *postgresql_storage.PostgresqlStorage
}

// effectiveDecls check provider decls and download provider binary file, return the effective providers
func effectiveDecls(ctx context.Context, decls []*config.ProviderDecl) (effects []*config.ProviderDecl, errlogs []string) {
	namespace, _, err := utils.Home()
	if err != nil {
		errlogs = append(errlogs, err.Error())
		return
	}
	provider := registry.NewProviderRegistry(namespace)
	ui.Successf("Selefra has been successfully installed providers!\n\n")
	ui.Successf("Checking Selefra provider updates......\n")

	for _, decl := range decls {
		configVersion := decl.Version
		prov := registry.Provider{
			Name:    decl.Name,
			Version: decl.Version,
			Source:  "",
			Path:    decl.Path,
		}
		pp, err := provider.Download(ctx, prov, true)
		if err != nil {
			ui.Errorf("%s@%s failed updated：%s", decl.Name, decl.Version, err.Error())
			errlogs = append(errlogs, err.Error())
			continue
		} else {
			decl.Path = pp.Filepath
			decl.Version = pp.Version
			err = tools.AppendProviderDecl(pp, nil, configVersion)
			if err != nil {
				ui.Errorf("%s@%s failed updated：%s", decl.Name, decl.Version, err.Error())
				errlogs = append(errlogs, err.Error())
				continue
			}
			effects = append(effects, decl)
			ui.Successf("	%s@%s all ready updated!\n", decl.Name, decl.Version)
		}
	}

	return effects, nil
}

func Sync(ctx context.Context) (lockSlice []lockStruct, err error) {
	// load and check config
	ui.Infof("Initializing provider plugins...\n\n")
	rootConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	if err = test.CheckSelefraConfig(ctx, rootConfig); err != nil {
		_ = httpClient.TrySetUpStage(global.RelvPrjName(), httpClient.Failed)
		return nil, err
	}

	if _, err := grpcClient.UploadLogStatus(); err != nil {
		ui.Errorln(err.Error())
	}

	var errored bool

	providerDecls, errLogs := effectiveDecls(ctx, rootConfig.Selefra.ProviderDecls)

	ui.Successf("Selefra has been finished update providers!\n")

	global.SetStage("pull")
	for _, decl := range providerDecls {
		prvds := tools.ProvidersByID(rootConfig, decl.Name)
		for _, prvd := range prvds {

			// build a postgresql storage
			schemaKey := config.GetSchemaKey(decl, *prvd)
			store, err := pgstorage.PgStorageWithMeta(ctx, &schema.ClientMeta{
				ClientLogger: logger.NewSchemaLoggeer(),
			}, pgstorage.WithSearchPath(config.GetSchemaKey(decl, *prvd)))
			if err != nil {
				errored = true
				ui.Errorf("%s@%s failed updated：%s", decl.Name, decl.Version, err.Error())
				errLogs = append(errLogs, fmt.Sprintf("%s@%s failed updated：%s", decl.Name, decl.Version, err.Error()))
				continue
			}

			// try lock
			// TODO: check unlock
			uuid := id_util.RandomId()
			for {
				err = store.Lock(ctx, schemaKey, uuid)
				if err == nil {
					lockSlice = append(lockSlice, lockStruct{
						SchemaKey: schemaKey,
						Uuid:      uuid,
						Storage:   store,
					})
					break
				}
				time.Sleep(5 * time.Second)
			}

			// check if cache expired
			expired, _ := tools.CacheExpired(ctx, store, prvd.Cache)
			if !expired {
				ui.Successf("%s %s@%s pull infrastructure data:\n", prvd.Name, decl.Name, decl.Version)
				ui.Print(fmt.Sprintf("Pulling %s@%s Please wait for resource information ...", decl.Name, decl.Version), false)
				ui.Successf("	%s@%s all ready use cache!\n", decl.Name, decl.Version)
				continue
			}

			// if expired, fetch new data
			err = fetch.Fetch(ctx, decl, prvd)
			if err != nil {
				ui.Errorf("%s %s Synchronization failed：%s", decl.Name, decl.Version, err.Error())
				errored = true
				continue
			}

			// set fetch time
			if err := pgstorage.SetStorageValue(ctx, store, config.GetCacheKey(), time.Now().Format(time.RFC3339)); err != nil {
				ui.Warningf("%s %s set cache time failed：%s", decl.Name, decl.Version, err.Error())
				errored = true
				continue
			}
		}
	}
	if errored {
		ui.Errorf(`
This may be exception, view detailed exception in %s .
`, filepath.Join(global.WorkSpace(), "logs"))
	}

	return lockSlice, nil
}
