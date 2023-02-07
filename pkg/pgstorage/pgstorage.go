package pgstorage

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-provider-sdk/storage"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/httpClient"
	"github.com/selefra/selefra/pkg/oci"
	"github.com/selefra/selefra/ui"
)

type Option func(pgopts *postgresql_storage.PostgresqlStorageOptions)

func DefaultPgStorageOpts() *postgresql_storage.PostgresqlStorageOptions {
	dsn := getDsn()

	pgopts := postgresql_storage.NewPostgresqlStorageOptions(dsn)

	return pgopts
}

func WithSearchPath(searchPath string) Option {
	return func(pgopts *postgresql_storage.PostgresqlStorageOptions) {
		pgopts.SearchPath = searchPath
	}
}

func PgStorage(ctx context.Context, opts ...Option) (*postgresql_storage.PostgresqlStorage, *schema.Diagnostics) {
	pgopts := DefaultPgStorageOpts()

	for _, opt := range opts {
		opt(pgopts)
	}

	return postgresql_storage.NewPostgresqlStorage(ctx, pgopts)
}

func Storage(ctx context.Context, opts ...Option) (storage.Storage, *schema.Diagnostics) {
	pgopts := DefaultPgStorageOpts()

	for _, opt := range opts {
		opt(pgopts)
	}

	return storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, pgopts)
}

func getDsn() (dsn string) {
	var err error
	if global.Token() != "" && global.RelvPrjName() != "" {
		dsn, err = httpClient.GetDsn(global.Token())
		if err != nil {
			ui.Errorln(err.Error())
			return ""
		}
	}

	err = oci.RunDB()
	if err != nil {
		ui.Errorln(err.Error())
		return ""
	}
	db := &config.DB{
		Driver:   "",
		Type:     "postgres",
		Username: "postgres",
		Password: "pass",
		Host:     "localhost",
		Port:     "15432",
		Database: "postgres",
		SSLMode:  "disable",
		Extras:   nil,
	}
	dsn = "host=" + db.Host + " user=" + db.Username + " password=" + db.Password + " port=" + db.Port + " dbname=" + db.Database + " " + "sslmode=disable"
	return
}
