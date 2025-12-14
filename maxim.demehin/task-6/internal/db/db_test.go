package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TvoyBatyA1234/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	selectAllNames = "SELECT name FROM users"
	selectUnique   = "SELECT DISTINCT name FROM users"
)

var TestError = errors.New("test error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Maria")
	mock.ExpectQuery(selectUnique).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alex", "Maria"}, names)
}

func TestGetUniqueNames_QueryFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(selectUnique).WillReturnError(TestError)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "db query")
	require.Nil(t, names)
}

func TestGetUniqueNames_InvalidData(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(selectUnique).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows scanning")
	require.Nil(t, names)
}

func TestGetUniqueNames_RowIssue(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Peter")
	rows.RowError(0, TestError)
	mock.ExpectQuery(selectUnique).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "rows error")
	require.Nil(t, names)
}
