package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

        "alina.duhanina/task-6/internal/db"
)

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Equal(t, "Alice", names[0])
	assert.Equal(t, "Bob", names[1])
}

func TestDBService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(sql.ErrConnDone)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.RowError(0, sql.ErrTxDone)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Len(t, names, 2)
	assert.Equal(t, "Alice", names[0])
	assert.Equal(t, "Bob", names[1])
}

func TestDBService_GetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(sql.ErrConnDone)

	service := db.New(mockDB)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.RowError(0, sql.ErrTxDone)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	service := db.New(mockDB)
	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	assert.Equal(t, mockDB, service.DB)
}
