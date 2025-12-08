package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	selectNamesQuery       = "SELECT name FROM users"
	selectUniqueNamesQuery = "SELECT DISTINCT name FROM users"
)

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Charlie")

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"Alice", "Bob", "Charlie"}, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Empty(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(selectNamesQuery).WillReturnError(sql.ErrConnDone)

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow(nil)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(1, sql.ErrTxDone)

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Alice")

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"Alice", "Bob", "Alice"}, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Empty(t, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnError(sql.ErrConnDone)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "db query")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestNew(t *testing.T) {
	t.Parallel()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	assert.NotNil(t, service)
	assert.Equal(t, db, service.DB)
}

func TestDBService_GetNames_WithCloseError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	mock.ExpectClose().WillReturnError(sql.ErrConnDone)

	service := New(db)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice"}, names)

	db.Close()
}

func TestDBService_GetNames_CloseError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
	rows = rows.CloseError(errors.New("close error"))

	mock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"Alice"}, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow(nil)

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows scanning")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, sql.ErrTxDone)

	mock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	service := New(db)
	names, err := service.GetUniqueNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "rows error")

	require.NoError(t, mock.ExpectationsWereMet())
}
