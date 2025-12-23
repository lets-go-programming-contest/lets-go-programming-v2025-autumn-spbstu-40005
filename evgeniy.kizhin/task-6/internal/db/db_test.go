package db_test

import (
	"errors"
	"testing"

	"evgeniy.kizhin/task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	errQueryFailed  = errors.New("query failed")
	errRowRead      = errors.New("row read error")
	errSelectFailed = errors.New("select failed")
)

func TestNewAndGetNames(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	svc := db.New(dbConn)

	t.Run("New returns service with same DB", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, dbConn, svc.DB)
	})

	t.Run("GetNames returns list on success", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow("Alice")
		rows.AddRow("Bob")
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := svc.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetNames returns error when query fails", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQueryFailed)

		names, err := svc.GetNames()
		require.ErrorContains(t, err, "db query")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetNames returns error scanning nil value", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow(nil)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := svc.GetNames()
		require.ErrorContains(t, err, "rows scanning")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetNames returns error when row has RowError", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow("Alice")
		rows.RowError(0, errRowRead)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := svc.GetNames()
		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = dbConn.Close() })

	svc := db.New(dbConn)

	t.Run("GetUniqueNames returns list on success", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow("Alice")
		rows.AddRow("Bob")
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := svc.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Alice", "Bob"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUniqueNames returns error when query fails", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errSelectFailed)

		names, err := svc.GetUniqueNames()
		require.ErrorContains(t, err, "db query")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUniqueNames scanning nil value returns error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow(nil)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := svc.GetUniqueNames()
		require.ErrorContains(t, err, "rows scanning")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetUniqueNames rows RowError returns error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})
		rows.AddRow("Alice")
		rows.RowError(0, errRowRead)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := svc.GetUniqueNames()
		require.ErrorContains(t, err, "rows error")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
