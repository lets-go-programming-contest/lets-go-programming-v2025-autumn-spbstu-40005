package db_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TvoyBatyA1234/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	selectAllNames = "SELECT name FROM users"
	selectUnique   = "SELECT DISTINCT name FROM users"
)

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
