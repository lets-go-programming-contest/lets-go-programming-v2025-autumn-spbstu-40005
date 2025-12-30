package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"nikita.brevnov/task-6/internal/db"
)

const (
	querySelectNames    = "SELECT name FROM users"
	querySelectDistinct = "SELECT DISTINCT name FROM users"
)

var ErrDB = errors.New("database error")

func TestFetchNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Mike")
	mock.ExpectQuery(querySelectNames).WillReturnRows(rows)

	result, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alex", "Mike"}, result)
}

func TestFetchNames_QueryFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(querySelectNames).WillReturnError(ErrDB)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "query")
	require.Nil(t, names)
}

func TestFetchNames_ScanFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(querySelectNames).WillReturnRows(rows)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan")
	require.Nil(t, names)
}

func TestFetchNames_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("John")
	rows.RowError(0, ErrDB)
	mock.ExpectQuery(querySelectNames).WillReturnRows(rows)

	names, err := service.GetNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "row")
	require.Nil(t, names)
}

func TestFetchUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Mike")
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alex", "Mike"}, names)
}

func TestFetchUniqueNames_QueryFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(querySelectDistinct).WillReturnError(ErrDB)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "query")
	require.Nil(t, names)
}

func TestFetchUniqueNames_ScanFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "scan")
	require.Nil(t, names)
}

func TestFetchUniqueNames_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("John")
	rows.RowError(0, ErrDB)
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Contains(t, err.Error(), "row")
	require.Nil(t, names)
}

func TestCreateService(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	svc := db.New(mockDB)
	require.Equal(t, mockDB, svc.DB)
}
