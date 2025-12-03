package db_test

import (
	"errors"
	"testing"

	"gordey.shapkov/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type rowTestDb struct {
	names       []string
	errExpected error
}

var testTable = []rowTestDb{
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names:       nil,
		errExpected: errors.New("empty names"),
	},
}

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}
	for _, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(mockDbRows(row.names)).WillReturnError(row.errExpected)

		names, err := dbService.GetNames()
		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error: %w, actual error: %w", row.errExpected, err)
			require.Nil(t, names, "names must be nil")
			continue
		}
		require.NoError(t, err, "error must be nil")
		require.Equal(t, row.names, names, "expected names: %s, actual names: %s", row.names, names)
	}

	row := rowTestDb{
		names:       []string{},
		errExpected: errors.New("err"),
	}

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(mockDbRows(row.names).AddRow(nil)).WillReturnError(row.errExpected)

	names, err := dbService.GetNames()
	if row.errExpected != nil {
		require.ErrorIs(t, err, row.errExpected, "expected error: %w, actual error: %w", row.errExpected, err)
		require.Nil(t, names, "names must be nil")
	} else {
		require.NoError(t, err, "error must be nil")
		require.Equal(t, row.names, names, "expected names: %s, actual names: %s", row.names, names)
	}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("network lost"))
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err = dbService.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)
}

func TestGetUniqueNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}
	for _, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(mockDbRows(row.names)).WillReturnError(row.errExpected)

		names, err := dbService.GetUniqueNames()
		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected, "row: %d, expected error: %w, actual error: %w", row.errExpected, err)
			require.Nil(t, names, "names must be nil")
			continue
		}
		require.NoError(t, err, "error must be nil")
		require.Equal(t, row.names, names, "expected names: %s, actual names: %s", row.names, names)
	}

	row := rowTestDb{
		names:       []string{},
		errExpected: errors.New("err"),
	}

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(mockDbRows(row.names).AddRow(nil)).WillReturnError(row.errExpected)

	names, err := dbService.GetUniqueNames()
	if row.errExpected != nil {
		require.ErrorIs(t, err, row.errExpected, "expected error: %w, actual error: %w", row.errExpected, err)
		require.Nil(t, names, "names must be nil")
	} else {
		require.NoError(t, err, "error must be nil")
		require.Equal(t, row.names, names, "expected names: %s, actual names: %s", row.names, names)
	}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("network lost"))
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err = dbService.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)
}

func mockDbRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})

	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}
