package main

import (
	"flag"
	"fmt"

	config "github.com/josemiguelmelo/godbmigrate/configuration"
	migrator "github.com/josemiguelmelo/godbmigrate/migrator"
	providers "github.com/josemiguelmelo/godbmigrate/providers"
)

const (
	defaultConfigName     = "godbmigrate.yml"
	defaultConfigLocation = "."

	defaultMigrationFolder = "./db/migrations"
)

func main() {
	runMigration := flag.Bool("migrate", false, "Run migrations")
	runMigrationRollback := flag.Bool("migrate-rollback", false, "Run migrations")

	configName := flag.String("config", defaultConfigName, "Configuration file name")
	configLocation := flag.String("config-location", defaultConfigLocation, "Configuration file location")

	migrationFolder := flag.String("migration-folder", defaultMigrationFolder, "Migration folders")
	flag.Parse()

	if *runMigration && *runMigrationRollback {
		panic("Cannot use migration and migration-rollback together")
	}
	if !*runMigration && !*runMigrationRollback {
		fmt.Println("Use one of these: -migrate or -migrate-rollback")
		return
	}

	runner := runner(*configName, *configLocation, *migrationFolder)

	if *runMigration {
		runner.RunUp()
	}

	if *runMigrationRollback {
		runner.RunDown()
		return
	}
}

func runner(configName, configLocation, migrationFolder string) *migrator.MigrationRunner {
	config, err := config.LoadConfiguration(configName, configLocation)
	if err != nil {
		panic(err)
	}

	provider := getDbProvider(config.Database.Provider, config.Database)
	return migrator.NewMigrationRunner(provider, migrationFolder)
}

func getDbProvider(provider string, dbConfig config.DatabaseConfiguration) migrator.MigratorInterface {
	switch provider {
	case "postgres":
		if dbConfig.ConnectionUri != "" {
			return providers.NewPostgresMigratorFromConnectionString(dbConfig.ConnectionUri)
		} else {
			return providers.NewPostgresMigrator(
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.Username,
				dbConfig.Password,
				dbConfig.Database,
				false,
			)
		}
	}

	panic("Provider not supported")
}
