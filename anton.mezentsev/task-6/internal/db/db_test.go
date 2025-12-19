package db_test

import (
	"errors"
	"testing"

	"anton.mezentsev/task-6/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

const (
	queryNames       = "SELECT name FROM users"
	queryUniqueNames = "SELECT DISTINCT name FROM users"
)

var ErrExpected = errors.New("expected error")

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228")
	mock.ExpectQuery(queryNames).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Ivan", "Gena228"}, names)
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}
	mock.ExpectQuery(queryNames).WillReturnError(ErrExpected)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "db query")
	require.Nil(t, names)
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Ivan").AddRow("Gena228")
	mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Ivan", "Gena228"}, names)
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}
	mock.ExpectQuery(queryUniqueNames).WillReturnError(ErrExpected)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "db query")
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
