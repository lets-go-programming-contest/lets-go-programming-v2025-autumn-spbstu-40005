package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/smirnov-vladislav/task-6/internal/db"
)

const (
	selectNamesQuery      = "SELECT name FROM users"
	selectDistinctQuery   = "SELECT DISTINCT name FROM users"
)

var errTest = errors.New("test error")

func closeMockDB(t *testing.T, mockDB *sql.DB) {
	t.Helper()
        if err := mockDB.Close(); err != nil {
        	t.Logf("Failed to close mock DB: %v", err)
        }
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Jane")
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"John", "Jane"}, names)
}

func TestGetNames_EmptyResult(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Empty(t, names)
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(selectNamesQuery).WillReturnError(errTest)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.ErrorContains(t, err, "db query")
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.ErrorContains(t, err, "rows scanning")
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")
	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Empty(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Test")
	rows.RowError(0, errTest)
	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
	require.ErrorContains(t, err, "rows error")
}

func TestNew(t *testing.T) {
	t.Parallel()
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer closeMockDB(t, mockDB)

	service := db.New(mockDB)
	require.NotNil(t, service)
	require.Equal(t, mockDB, service.DB)
}
