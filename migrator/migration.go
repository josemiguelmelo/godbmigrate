package migrator

import (
	"os"
	"strings"
)

// Migration Migration struct
type Migration struct {
	migrationFileLocation string
	name                  string
}

// NewMigration Instantiate new migration
func NewMigration(migrationLocation string) *Migration {
	return &Migration{
		migrationFileLocation: migrationLocation,
		name:                  migrationName(migrationLocation),
	}
}

// Up Get Up migration
func (m *Migration) Up() (string, error) {
	migrationContent, err := m.content()
	if err != nil {
		return "", err
	}
	upAndDownArray := strings.Split(migrationContent, downSection)

	return strings.Replace(upAndDownArray[0], upSection, "", 1), nil
}

// Down Get down migration
func (m *Migration) Down() (string, error) {
	migrationContent, err := m.content()
	if err != nil {
		return "", err
	}

	upAndDownArray := strings.Split(migrationContent, downSection)

	return strings.Replace(upAndDownArray[1], downSection, "", 1), nil
}

func (m *Migration) content() (string, error) {
	reader := NewMigrationReader(m.migrationFileLocation)
	return reader.Read()
}

func migrationName(migrationLocation string) string {
	path := strings.Split(migrationLocation, string(os.PathSeparator))
	return path[len(path)-1]
}
