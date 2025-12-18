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
)

type nameTestCase struct {
	name        string
	returnErr   error
	returnNames []string
	wantErr     bool
}

var nameTestTable = []nameTestCase{
	{
		name:        "success",
		returnNames: []string{"Polina", "Artemiy"},
	},
	{
		name:      "query error",
		returnErr: ErrDBConnectionFailed,
		wantErr:   true,
	},
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	for _, tt := range nameTestTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := myDb.New(mockDB)

			if tt.returnErr != nil {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(tt.returnErr)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range tt.returnNames {
					rows = rows.AddRow(name)
				}
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			}

			names, err := service.GetNames()

			if tt.wantErr {
				require.ErrorContains(t, err, "db query:")
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.returnNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("scan error - nil value", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetNames()
		require.ErrorContains(t, err, "rows scanning:")
		require.Nil(t, names)
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
		require.ErrorContains(t, err, "rows error:")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	for _, tt := range nameTestTable {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := myDb.New(mockDB)

			if tt.returnErr != nil {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(tt.returnErr)
			} else {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range tt.returnNames {
					rows = rows.AddRow(name)
				}
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			}

			names, err := service.GetUniqueNames()

			if tt.wantErr {
				require.ErrorContains(t, err, "db query:")
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.returnNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("scan error - nil value", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		service := myDb.New(mockDB)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		names, err := service.GetUniqueNames()
		require.ErrorContains(t, err, "rows scanning:")
		require.Nil(t, names)
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
		require.ErrorContains(t, err, "rows error:")
		require.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
