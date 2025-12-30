package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"kirill.kharlamov/task-6/internal/db"
)

const (
	queryNames       = "SELECT name FROM users"
	queryUniqueNames = "SELECT DISTINCT name FROM users"
)

var ErrTest = errors.New("test error")

func TestGetNames_Standard(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Maria")
	mock.ExpectQuery(queryNames).WillReturnRows(rows)

	result, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alex", "Maria"}, result)
}

func TestGetNames_QueryFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(queryNames).WillReturnError(ErrTest)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "db query")
	require.Nil(t, names)
}

func TestGetNames_ScanFailure(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(queryNames).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)
}

func TestGetNames_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Dmitry")
	rows.RowError(0, ErrTest)
	mock.ExpectQuery(queryNames).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
}

func TestGetUniqueNames_Standard(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Unique1").AddRow("Unique2")
	mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Unique1", "Unique2"}, names)
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(queryUniqueNames).WillReturnError(ErrTest)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "db query")
	require.Nil(t, names)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Sergey")
	rows.RowError(0, ErrTest)
	mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}
