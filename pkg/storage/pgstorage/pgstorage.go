package pgstorage

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/pkg/modules/module"
	"strings"
)

type Option func(pgopts *postgresql_storage.PostgresqlStorageOptions)

//func DefaultPgStorageOpts() *postgresql_storage.PostgresqlStorageOptions {
//	dsn := getDsn()
//
//	pgopts := postgresql_storage.NewPostgresqlStorageOptions(dsn)
//
//	return pgopts
//}

func WithSearchPath(searchPath string) Option {
	return func(pgopts *postgresql_storage.PostgresqlStorageOptions) {
		pgopts.SearchPath = searchPath
	}
}

//func PgStorageWithMeta(ctx context.Context, meta *schema.ClientMeta, opts ...Option) (*postgresql_storage.PostgresqlStorage, error) {
//	pgopts := DefaultPgStorageOpts()
//
//	for _, opt := range opts {
//		opt(pgopts)
//	}
//
//	storage, diag := postgresql_storage.NewPostgresqlStorage(ctx, pgopts)
//	if diag != nil {
//		if diag != nil {
//			err := cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
//			if err != nil {
//				return nil, errors.New(`The database maybe not ready.
//		You can execute the following command to install the official database image.
//		docker run --name selefra_postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres\n`)
//			}
//		}
//	}
//
//	storage.SetClientMeta(meta)
//
//	return storage, nil
//}
//
//func PgStorage(ctx context.Context, opts ...Option) (*postgresql_storage.PostgresqlStorage, *schema.Diagnostics) {
//	pgopts := DefaultPgStorageOpts()
//
//	for _, opt := range opts {
//		opt(pgopts)
//	}
//
//	return postgresql_storage.NewPostgresqlStorage(ctx, pgopts)
//}
//
//func Storage(ctx context.Context, opts ...Option) (storage.Storage, *schema.Diagnostics) {
//	pgopts := DefaultPgStorageOpts()
//
//	for _, opt := range opts {
//		opt(pgopts)
//	}
//
//	return storage_factory.NewStorage(ctx, storage_factory.StorageTypePostgresql, pgopts)
//}

func GetStorageValue(ctx context.Context, storage *postgresql_storage.PostgresqlStorage, key string) (string, error) {
	v, diag := storage.GetValue(ctx, key)
	if diag != nil {
		err := cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		if err != nil {
			return "", err
		}
	}
	return v, nil
}

func SetStorageValue(ctx context.Context, storage *postgresql_storage.PostgresqlStorage, key, value string) error {
	if diag := storage.SetKey(ctx, key, value); diag != nil {
		err := cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		if err != nil {
			return err
		}
	}

	return nil
}

//func getDsn() (dsn string) {
//	var err error
//	if global.Token() != "" && global.RelvPrjName() != "" {
//		dsn, err = http_client.GetDsn(global.Token())
//		if err != nil {
//			ui.Errorln(err.Error())
//			return ""
//		}
//	}
//
//	err = oci.RunDB()
//	if err != nil {
//		ui.Errorln(err.Error())
//		return ""
//	}
//	db := &config.DB{
//		Driver:   "",
//		Type:     "postgres",
//		Username: "postgres",
//		Password: "pass",
//		Host:     "localhost",
//		Port:     "15432",
//		Database: "postgres",
//		SSLMode:  "disable",
//		Extras:   nil,
//	}
//	dsn = "host=" + db.Host + " user=" + db.Username + " password=" + db.Password + " port=" + db.Port + " dbname=" + db.Database + " " + "sslmode=disable"
//	return
//}

// GetSchemaKey return provider schema named <required.name>_<required_version>_<provider_name>
func GetSchemaKey(providerName, providerVersion string, providerConfigurationBlock *module.ProviderBlock) string {
	sourceArr := strings.Split(providerName, "/")
	var source string
	if len(sourceArr) > 1 {
		source = strings.Replace(sourceArr[1]+"@"+providerVersion, "/", "_", -1)
	} else {
		source = strings.Replace(sourceArr[0]+"@"+providerVersion, "/", "_", -1)
	}
	source = strings.Replace(source, "@", "_", -1)
	source = strings.Replace(source, ".", "", -1)
	if providerConfigurationBlock != nil {
		source = source + "_" + providerConfigurationBlock.Name
	}
	return source
}
