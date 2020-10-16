package migrator

import (
	"fmt"

	"github.com/google/uuid"
)

type MigrationRunner struct {
	migrator         MigratorInterface
	reader           *MigrationReader
	migrationsFolder string
}

func NewMigrationRunner(migrator MigratorInterface, location string) *MigrationRunner {
	migrationReader := NewMigrationReaderStatic()
	return &MigrationRunner{
		migrator:         migrator,
		reader:           migrationReader,
		migrationsFolder: location,
	}
}

func (m *MigrationRunner) migrationSuccessfullyRun(migration Migration, runID string) {
	fmt.Printf("%s run successfully!\n", migration.name)
	m.migrator.InsertIntoMigrationsTable(migration.name, runID)
}

func (m *MigrationRunner) migrationSuccessfullyRolledback(migration Migration, runID string) {
	fmt.Printf("%s rolled back successfully!\n", migration.name)
	m.migrator.DeleteFromMigrationsTable(migration.name, runID)
}

func (m *MigrationRunner) RunDown() {
	migrationRunID, err := m.migrator.LastMigrationRunID()
	if err != nil {
		panic("No previous migration found." + err.Error())
	}

	migrations, err := m.migrator.FetchMigrationsByRunID(migrationRunID)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch last migration run identifier: %s", err.Error()))
	}
	allMigrations := m.reader.ListAllMigrations(m.migrationsFolder)

	m.migrator.StartTransaction()

	for _, migrationName := range migrations {
		for _, migration := range allMigrations {
			if migration.name == migrationName {
				fmt.Printf("=== Rolling migration %s ===\n", migration.name)

				migrationQuery, err := migration.Down()
				if err != nil {
					fmt.Printf("Error parsing migration %s: %s\n", migration.name, err)
					m.migrator.RollbackTransaction()
					panic(err)
				}

				err = m.migrator.RunMigration(migrationQuery)
				if err != nil {
					fmt.Printf("Error rolling back migration %s: %s\n", migration.name, err)
					m.migrator.RollbackTransaction()
					panic(err)
				} else {
					m.migrationSuccessfullyRolledback(migration, migrationRunID)
				}

			}
		}
	}

	m.migrator.CommitTransaction()
}

func (m *MigrationRunner) RunUp() {
	m.migrator.CreateMigrationsTable()

	migrationUUUID, _ := uuid.NewRandom()
	migrationID := migrationUUUID.String()

	migrations := m.reader.ListAllMigrations(m.migrationsFolder)

	m.migrator.StartTransaction()
	for _, migration := range migrations {

		if m.migrator.MigrationAlreadyRun(migration.name) {
			continue
		}

		fmt.Printf("=== Running migration %s ===\n", migration.name)

		migrationQuery, err := migration.Up()
		if err != nil {
			fmt.Printf("Error parsing migration %s: %s\n", migration.name, err)
			m.migrator.RollbackTransaction()
			panic(err)
		}
		err = m.migrator.RunMigration(migrationQuery)

		if err != nil {
			fmt.Printf("Error running migration %s: %s\n", migration.name, err)
			m.migrator.RollbackTransaction()
			panic(err)
		} else {
			m.migrationSuccessfullyRun(migration, migrationID)
		}

		fmt.Println("")
	}

	m.migrator.CommitTransaction()
}
