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

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errContains string
		wantNames   []string
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Polina").
					AddRow("Artemiy")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Polina", "Artemiy"},
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnError(ErrDBConnectionFailed)
			},
			wantErr:     true,
			errContains: "db query:",
		},
		{
			name: "scan error - nil value",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			wantErr:     true,
			errContains: "rows scanning:",
		},
		{
			name: "rows.Err after scan",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
				rows.RowError(0, ErrNetworkAfterScan)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "rows error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := myDb.New(mockDB)

			tt.mockSetup(mock)

			names, err := service.GetNames()

			if tt.wantErr {
				require.ErrorContains(t, err, tt.errContains)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mockSetup   func(mock sqlmock.Sqlmock)
		wantErr     bool
		errContains string
		wantNames   []string
	}{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Polina").
					AddRow("Artemiy")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantNames: []string{"Polina", "Artemiy"},
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(ErrDBConnectionFailed)
			},
			wantErr:     true,
			errContains: "db query:",
		},
		{
			name: "scan error - nil value",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			wantErr:     true,
			errContains: "rows scanning:",
		},
		{
			name: "rows.Err after scan",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Polina")
				rows.RowError(0, ErrNetworkAfterScan)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr:     true,
			errContains: "rows error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := myDb.New(mockDB)

			tt.mockSetup(mock)

			names, err := service.GetUniqueNames()

			if tt.wantErr {
				require.ErrorContains(t, err, tt.errContains)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
