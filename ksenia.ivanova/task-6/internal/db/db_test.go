package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"task-6/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	errQueryFailed  = errors.New("query failed")
	errRowRead      = errors.New("row read error")
	errSelectFailed = errors.New("select failed")
)

func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, db.DBService) {
	t.Helper()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = dbConn.Close()
	})

	return dbConn, mock, db.New(dbConn)
}

func TestNewServiceReturnsSameDB(t *testing.T) {
	dbConn, _, svc := setupDB(t)
	t.Parallel()
	require.Equal(t, dbConn, svc.DB)
}

func TestGetNamesSuccess(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.AddRow("Bob")
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := svc.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesQueryFailed(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQueryFailed)

	names, err := svc.GetNames()
	require.ErrorContains(t, err, "db query")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesScanningNilValue(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow(nil)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := svc.GetNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesRowError(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.RowError(0, errRowRead)
	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := svc.GetNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesSuccess(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.AddRow("Bob")
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := svc.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesQueryFailed(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errSelectFailed)

	names, err := svc.GetUniqueNames()
	require.ErrorContains(t, err, "db query")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesScanningNilValue(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow(nil)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := svc.GetUniqueNames()
	require.ErrorContains(t, err, "rows scanning")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesRowError(t *testing.T) {
	_, mock, svc := setupDB(t)
	t.Parallel()

	rows := sqlmock.NewRows([]string{"name"})
	rows.AddRow("Alice")
	rows.RowError(0, errRowRead)
	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := svc.GetUniqueNames()
	require.ErrorContains(t, err, "rows error")
	require.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}
