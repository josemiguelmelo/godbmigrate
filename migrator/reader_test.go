package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrationSorter(t *testing.T) {
	assert := assert.New(t)

	m := &MigrationReader{}
	migrationsList := []Migration{
		{
			name: "1_create_table_1",
		},
		{
			name: "3_create_table_3",
		},
		{
			name: "10_create_table_40",
		},
		{
			name: "2_create_table_2",
		},
	}

	assert.Equal("1_create_table_1", migrationsList[0].name)
	assert.Equal("3_create_table_3", migrationsList[1].name)
	assert.Equal("10_create_table_40", migrationsList[2].name)
	assert.Equal("2_create_table_2", migrationsList[3].name)

	m.sortMigrationsList(migrationsList)

	assert.Equal("1_create_table_1", migrationsList[0].name)
	assert.Equal("2_create_table_2", migrationsList[1].name)
	assert.Equal("3_create_table_3", migrationsList[2].name)
	assert.Equal("10_create_table_40", migrationsList[3].name)

}
