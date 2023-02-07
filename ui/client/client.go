package client

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/selefra/selefra-provider-sdk/storage"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/pgstorage"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/ui"
)

// TODO: Deprecated
type Client struct {
	//downloadProgress ui.Progress
	cfg           *config.SelefraConfig
	Providers     registry.Providers
	Registry      interface{}
	PluginManager interface{}
	Storage       storage.Storage
	instanceId    uuid.UUID
}

// TODO: Deprecated
func CreateClientFromConfig(ctx context.Context, cfg *config.SelefraConfig, instanceId uuid.UUID, provider *config.ProviderRequired, cp config.ProviderConfig) (*Client, error) {

	hub := new(interface{})
	pm := new(interface{})

	c := &Client{
		Storage:       nil,
		cfg:           cfg,
		Registry:      hub,
		PluginManager: pm,
		instanceId:    instanceId,
	}

	schema := config.GetSchemaKey(provider, cp)
	sto, diag := pgstorage.Storage(ctx, pgstorage.WithSearchPath(schema))
	if diag != nil {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		if err != nil {
			return nil, errors.New("failed to create pgstorage")
		}
	}
	if sto != nil {
		c.Storage = sto
	}

	//if cfg.GetDSN() != "" {
	//	options := postgres.NewPostgresqlStorageOptions(cfg.GetDSN())
	//	schema := config.GetSchemaKey(provider, cp)
	//	options.SearchPath = schema
	//	sto, diag := storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, options)
	//	if diag != nil {
	//		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
	//		if err != nil {
	//			return nil, errors.New("failed to create pgstorage")
	//		}
	//	}
	//	c.Storage = sto
	//}
	c.Providers = registry.Providers{}
	for _, rp := range cfg.Providers {
		c.Providers.Set(registry.Provider{Name: rp.Name, Version: rp.Version})
	}

	return c, nil
}
