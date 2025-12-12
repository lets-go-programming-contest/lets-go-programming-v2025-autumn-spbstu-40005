package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Polina").AddRow("Artemiy"))

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Polina", "Artemiy"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("db connection failed"))

		names, err := service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
		rows.RowError(0, errors.New("network error after scan"))
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - type mismatch", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

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
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Polina").AddRow("Artemiy"))

		names, err := service.GetUniqueNames()
		require.NoError(t, err)
		require.Equal(t, []string{"Polina", "Artemiy"}, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("query failed"))

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
		rows.RowError(0, errors.New("network error after scan"))
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error:")

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error - type mismatch", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetUniqueNames()
		require.Error(t, err)
		require.Contains(t, err.Error(), "rows scanning:")
		require.Nil(t, names)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
