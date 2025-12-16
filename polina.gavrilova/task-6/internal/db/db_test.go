package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	myDb "polina.gavrilova/task-6/internal/db"
)

var (
	ErrDBConnectionFailed = errors.New("db connection failed")
	ErrNetworkAfterScan   = errors.New("network error after scan")
	ErrQueryFailed        = errors.New("query failed")
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Polina").AddRow("Artemiy"))

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Polina", "Artemiy"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(ErrDBConnectionFailed)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
		rows.RowError(0, ErrNetworkAfterScan)
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - type mismatch", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Polina").AddRow("Artemiy"))

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Polina", "Artemiy"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(ErrQueryFailed)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
		rows.RowError(0, ErrNetworkAfterScan)
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - type mismatch", func(t *testing.T) {
		t.Parallel()
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
