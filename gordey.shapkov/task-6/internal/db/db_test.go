package db_test

import (
	"errors"
	"testing"

	"gordey.shapkov/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var testTable = []struct { //nolint:gochecknoglobals
	names       []string
	errWrap     string
	errExpected error
}{
	{
		names: []string{"Ivan", "Gena228"},
	},
	{
		names:       nil,
		errWrap:     "db query: ",
		errExpected: ErrExpected,
	},
}

var ErrExpected = errors.New("expected error")

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for _, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(helperMockDBRows(t, row.names)).
			WillReturnError(row.errExpected)

		names, err := dbService.GetNames()
		if row.errExpected != nil {
			require.ErrorContains(t, err, row.errWrap)
			require.Nil(t, names)

			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.names, names)
	}

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := dbService.GetNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, ErrExpected)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err = dbService.GetNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for _, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(helperMockDBRows(t, row.names)).
			WillReturnError(row.errExpected)

		names, err := dbService.GetUniqueNames()
		if row.errExpected != nil {
			require.ErrorContains(t, err, row.errWrap)
			require.Nil(t, names)

			continue
		}

		require.NoError(t, err)
		require.Equal(t, row.names, names)
	}

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	names, err := dbService.GetUniqueNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, ErrExpected)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err = dbService.GetUniqueNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
}

func helperMockDBRows(t *testing.T, names []string) *sqlmock.Rows {
	t.Helper()

	rows := sqlmock.NewRows([]string{"name"})

	for _, name := range names {
		rows = rows.AddRow(name)
	}

	return rows
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	require.Equal(t, mockDB, service.DB)
}
