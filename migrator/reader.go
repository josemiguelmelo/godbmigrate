package migrator

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	upSection   = "--- up! ---"
	downSection = "--- down! ---"
)

type MigrationReader struct {
	MigrationLocation *string
}

func NewMigrationReader(migrationLocation string) *MigrationReader {
	return &MigrationReader{
		MigrationLocation: &migrationLocation,
	}
}

func NewMigrationReaderStatic() *MigrationReader {
	return &MigrationReader{}
}

func (m *MigrationReader) ListAllMigrations(migrationRootFolder string) []Migration {
	var files []Migration

	err := filepath.Walk(migrationRootFolder, func(path string, info os.FileInfo, err error) error {
		if path != migrationRootFolder {
			files = append(files, *NewMigration(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func (m *MigrationReader) Read() (string, error) {
	if m.MigrationLocation == nil {
		return "", errors.New("Migration location is not set")
	}

	dataBytes, err := ioutil.ReadFile(*m.MigrationLocation)
	if err != nil {
		return "nil", err
	}

	return string(dataBytes), nil
}
