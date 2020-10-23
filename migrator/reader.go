package migrator

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	upSection   = "--- up! ---"
	downSection = "--- down! ---"
)

type MigrationReader struct {
	MigrationLocation *string
}

func NewMigrationReader(migrationLocation string) *MigrationReader {
	migrationLocation = normalizePath(migrationLocation)

	return &MigrationReader{
		MigrationLocation: &migrationLocation,
	}
}

func NewMigrationReaderStatic() *MigrationReader {
	return &MigrationReader{}
}

func normalizePath(path string) string {
	return strings.ReplaceAll(path, "/", string(os.PathSeparator))
}

func migrationIndexFromName(name string) (int, error) {
	nameParts := strings.Split(name, "_")
	return strconv.Atoi(nameParts[0])
}

func (m *MigrationReader) sortMigrationsList(migrations []Migration) []Migration {
	sort.SliceStable(migrations, func(i, j int) bool {
		firstNameIndex, err := migrationIndexFromName(migrations[i].name)
		if err != nil {
			panic(err)
		}

		secondNameIndex, err := migrationIndexFromName(migrations[j].name)
		if err != nil {
			panic(err)
		}

		return firstNameIndex < secondNameIndex
	})
	return migrations
}

func (m *MigrationReader) ListAllMigrations(migrationRootFolder string) []Migration {
	var migrations []Migration

	migrationRootFolder = normalizePath(migrationRootFolder)

	err := filepath.Walk(migrationRootFolder, func(path string, info os.FileInfo, err error) error {
		if path != migrationRootFolder {
			migrations = append(migrations, *NewMigration(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return m.sortMigrationsList(migrations)
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
