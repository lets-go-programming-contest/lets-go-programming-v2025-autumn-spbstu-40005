package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"nikita.brevnov/task-6/internal/db"
)

type testCase struct {
	names []string
	err   error
}

var testScenarios = []testCase{
	{
		names: []string{"Иван", "Петр"},
	},
	{
		names: nil,
		err:   errors.New("DBError"),
	},
}

func TestFetchUserNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}

	service := db.DBService{DB: mockDB}
	for i, scenario := range testScenarios {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(buildRows(scenario.names)).
			WillReturnError(scenario.err)
		result, err := service.GetNames()

		if scenario.err != nil {
			require.Error(t, err, "row: %d, expected error", i)
			require.Nil(t, result, "row: %d, result should be nil", i)
			continue
		}
		require.NoError(t, err, "row: %d, error should be nil", i)
		require.Equal(t, scenario.names, result, "row: %d, names don't match", i)
	}
}

func buildRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows.AddRow(name)
	}
	return rows
}
