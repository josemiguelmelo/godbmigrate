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
	return &MigrationReader{
		MigrationLocation: &migrationLocation,
	}
}

func NewMigrationReaderStatic() *MigrationReader {
	return &MigrationReader{}
}

func (m *MigrationReader) sortMigrationsList(migrations []Migration) {
	sort.SliceStable(migrations, func(i, j int) bool {
		firstNameSplitted := strings.Split(migrations[i].name, "_")
		secondNameSplitted := strings.Split(migrations[j].name, "_")

		firstNameIndex, err := strconv.Atoi(firstNameSplitted[0])
		if err != nil {
			panic(err)
		}
		secondNameIndex, err := strconv.Atoi(secondNameSplitted[0])
		if err != nil {
			panic(err)
		}
		return firstNameIndex < secondNameIndex
	})
}

func (m *MigrationReader) ListAllMigrations(migrationRootFolder string) []Migration {
	var migrations []Migration

	err := filepath.Walk(migrationRootFolder, func(path string, info os.FileInfo, err error) error {
		if path != migrationRootFolder {
			migrations = append(migrations, *NewMigration(path))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return migrations
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
