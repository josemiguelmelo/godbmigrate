package migrator

type MigratorInterface interface {
	RunMigration(migrationQuery string) error
	StartTransaction()
	RollbackTransaction()
	CreateMigrationsTable() error
	InsertIntoMigrationsTable(migrationName string, runID string) error
	DeleteFromMigrationsTable(migrationName string, runID string) error
	MigrationAlreadyRun(migrationName string) bool
	CommitTransaction()
	LastMigrationRunID() (string, error)
	FetchMigrationsByRunID(runID string) ([]string, error)
}
