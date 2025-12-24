package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/smirnov-vladislav/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	selectNamesQuery    = "SELECT name FROM users"
	selectDistinctQuery = "SELECT DISTINCT name FROM users"
	msgDbQuery = "db query"
	msgRowsError = "rows error"
	msgRowsScan = "rows scanning"
)

var errTest = errors.New("test error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.New(mockDB)

	require.Equal(t, mockDB, service.DB)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Jane")

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	require.Equal(t, []string{"John", "Jane"}, names)
}

func TestGetNames_RowProblem(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Carl")
	rows.RowError(0, errTest)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.ErrorContains(t, err, msgRowsError)
	require.Nil(t, names)
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(selectNamesQuery).WillReturnError(errTest)

	names, err := service.GetNames()

	require.ErrorContains(t, err, msgDbQuery)
	require.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	names, err := service.GetNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, msgRowsScan)
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetUniqueNames_QueryFail(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(selectDistinctQuery).WillReturnError(errTest)

	names, err := service.GetUniqueNames()

	require.ErrorContains(t, err, msgDbQuery)
	require.Nil(t, names)
}

func TestGetUniqueNames_InvalData(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.ErrorContains(t, err, msgRowsScan)
	require.Nil(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()

	require.NoError(t, err)

	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Carl")

	rows.RowError(0, errTest)
	mock.ExpectQuery(selectDistinctQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Nil(t, names)
	require.ErrorContains(t, err, msgRowsError)
}
