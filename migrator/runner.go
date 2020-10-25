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

// NewMigrationRunner Instantiate new migration runner
func NewMigrationRunner(migrator MigratorInterface, location string) *MigrationRunner {
	migrationReader := NewMigrationReaderStatic()
	return &MigrationRunner{
		migrator:         migrator,
		reader:           migrationReader,
		migrationsFolder: location,
	}
}

func (m *MigrationRunner) migrationSuccessfullyRun(migration Migration, runID string) error {
	fmt.Printf("%s run successfully!\n", migration.name)
	return m.migrator.InsertIntoMigrationsTable(migration.name, runID)
}

func (m *MigrationRunner) migrationSuccessfullyRolledback(migration Migration, runID string) error {
	fmt.Printf("%s rolled back successfully!\n", migration.name)
	return m.migrator.DeleteFromMigrationsTable(migration.name, runID)
}

func (m *MigrationRunner) migrationFailed(err error) {
	m.migrator.RollbackTransaction()
	panic(err)
}

// RunDown Run rollback migrations
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
					m.migrationFailed(err)
				} else {
					if err = m.migrationSuccessfullyRolledback(migration, migrationRunID); err != nil {
						m.migrationFailed(err)
					}
				}
			}
		}
	}

	m.migrator.CommitTransaction()
}

// RunUp Run migration
func (m *MigrationRunner) RunUp() {
	err := m.migrator.CreateMigrationsTable()
	if err != nil {
		panic("Could not create migrations table" + err.Error())
	}

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
			m.migrationFailed(err)
		}
		err = m.migrator.RunMigration(migrationQuery)

		if err != nil {
			fmt.Printf("Error running migration %s: %s\n", migration.name, err)
			m.migrationFailed(err)
		} else {
			if err = m.migrationSuccessfullyRun(migration, migrationID); err != nil {
				m.migrationFailed(err)
			}
		}

		fmt.Println("")
	}

	m.migrator.CommitTransaction()
}
