package internaldb_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	internaldb "sergey.dribas/task-6/internal/db"
)

var (
	errRow   = errors.New("row error")
	errQuery = errors.New("query failed")
)

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expectNames []string
		expectErr   string
	}{
		{
			name: "success with multiple names",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Charlie")
				mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
			},
			expectNames: []string{"Alice", "Bob", "Charlie"},
			expectErr:   "",
		},
		{
			name: "success with empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
			},
			expectNames: nil,
			expectErr:   "",
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnError(errQuery)
			},
			expectNames: nil,
			expectErr:   "db query: query failed",
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{})
				rows.AddRow()
				mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
			},
			expectNames: nil,
			expectErr:   "rows scanning:",
		},
		{
			name: "rows.Err error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				rows = rows.RowError(0, errRow)
				mock.ExpectQuery(`^SELECT name FROM users$`).WillReturnRows(rows)
			},
			expectNames: nil,
			expectErr:   "rows error: row error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer mockDB.Close()

			service := internaldb.New(mockDB)

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			names, err := service.GetNames()

			if tt.expectErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		setupMock    func(sqlmock.Sqlmock)
		expectValues []string
		expectErr    string
	}{
		{
			name: "success with duplicates",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					AddRow("Alice")
				mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
			},
			expectValues: []string{"Alice", "Bob", "Alice"},
			expectErr:    "",
		},
		{
			name: "empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
			},
			expectValues: nil,
			expectErr:    "",
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnError(errQuery)
			},
			expectValues: nil,
			expectErr:    "db query: query failed",
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{})
				rows.AddRow()
				mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
			},
			expectValues: nil,
			expectErr:    "rows scanning:",
		},
		{
			name: "rows.Err error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Alice")
				rows = rows.RowError(0, errRow)
				mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).WillReturnRows(rows)
			},
			expectValues: nil,
			expectErr:    "rows error: row error",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := internaldb.New(mockDB)

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			values, err := service.GetUniqueNames()

			if tt.expectErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectErr)
				require.Nil(t, values)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectValues, values)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
