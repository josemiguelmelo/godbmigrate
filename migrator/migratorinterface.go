package migrator

type MigratorInterface interface {
	RunMigration(migrationQuery string) error
	StartTransaction()
	RollbackTransaction()
	CreateMigrationsTable() error
	InsertIntoMigrationsTable(migrationName string, runID string) error
	MigrationAlreadyRun(migrationName string) bool
	CommitTransaction()
}
