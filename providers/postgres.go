package providers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DefaultLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logger.Info, // Log level
		Colorful:      false,       // Disable color
	},
)

type PostgresMigrator struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewPostgresMigrator(host string, port int, username string, password string, database string, sslEnabled bool) *PostgresMigrator {
	migrator := new(PostgresMigrator)
	migrator.connect(host, port, username, password, database, sslEnabled)
	return migrator
}

func NewPostgresMigratorFromConnectionString(connectionString string) *PostgresMigrator {
	migrator := new(PostgresMigrator)
	migrator.connectWithConnectionString(connectionString)
	return migrator
}

func (m *PostgresMigrator) StartTransaction() {
	m.tx = m.db.Begin()
}

func (m *PostgresMigrator) RollbackTransaction() {
	m.tx.Rollback()
}

func (m *PostgresMigrator) CommitTransaction() {
	m.tx.Commit()
}

func (m *PostgresMigrator) CreateMigrationsTable() error {
	res := m.db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (id BIGSERIAL PRIMARY KEY NOT NULL, name VARCHAR(255) NOT NULL UNIQUE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), migration_id VARCHAR(255) NOT NULL)")
	return res.Error
}

func (m *PostgresMigrator) InsertIntoMigrationsTable(migrationName string, runID string) error {
	res := m.db.Exec("INSERT INTO schema_migrations(name, migration_id) VALUES (?, ?)", migrationName, runID)
	return res.Error
}

func (m *PostgresMigrator) MigrationAlreadyRun(migrationName string) bool {
	var res int
	m.db.Raw("SELECT count(*) as count from schema_migrations where name = ?", migrationName).Scan(&res)
	return res == 1
}

func (m *PostgresMigrator) RunMigration(migrationQuery string) error {
	return m.runMigration(migrationQuery)
}

func (m *PostgresMigrator) runMigration(query string) error {
	if strings.TrimSpace(query) != "" {
		tx := m.db.Exec(query)
		if tx.Error != nil {
			return tx.Error
		} else {
			return nil
		}
	} else {
		return errors.New("Migration is empty")
	}
}

func (m *PostgresMigrator) sslMode(sslEnabled bool) string {
	if sslEnabled {
		return "enable"
	}
	return "disable"
}

func (m *PostgresMigrator) connect(host string, port int, username string, password string, database string, sslEnabled bool) {
	dsn := fmt.Sprintf("user=%s password=%s DB.name=%s port=%d host=%s sslmode=%s", username, password, database, port, host, m.sslMode(sslEnabled))
	m.connectWithConnectionString(dsn)
}

func (m *PostgresMigrator) connectWithConnectionString(connectionString string) {
	var err error
	m.db, err = gorm.Open(
		postgres.Open(connectionString),
		&gorm.Config{Logger: DefaultLogger},
	)

	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, _ := m.db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
