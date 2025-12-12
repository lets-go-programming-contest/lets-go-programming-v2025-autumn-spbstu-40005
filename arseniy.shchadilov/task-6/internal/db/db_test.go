package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/arseniy.shchadilov/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	queryNames       = "SELECT name FROM users"
	queryUniqueNames = "SELECT DISTINCT name FROM users"
)

var (
	errExpected = errors.New("expected error")
)

type testCase struct {
	name          string
	setupMock     func(mock sqlmock.Sqlmock) *sqlmock.Rows
	expectedNames []string
	expectedErr   error
	errContains   string
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name: "success with multiple names",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("John")
				mock.ExpectQuery(queryNames).WillReturnRows(rows)

				return rows
			},
			expectedNames: []string{"Ivan", "John"},
			expectedErr:   nil,
		},
		{
			name: "success with empty result",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(queryNames).WillReturnRows(rows)

				return rows
			},
			expectedNames: []string{},
			expectedErr:   nil,
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				mock.ExpectQuery(queryNames).WillReturnError(errExpected)

				return nil
			},
			expectedNames: nil,
			expectedErr:   errExpected,
			errContains:   "db query",
		},
		{
			name: "scan error (nil value)",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(queryNames).WillReturnRows(rows)

				return rows
			},
			expectedNames: nil,
			expectedErr:   nil,
			errContains:   "rows scanning",
		},
		{
			name: "rows error",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					RowError(0, errExpected)
				mock.ExpectQuery(queryNames).WillReturnRows(rows)

				return rows
			},
			expectedNames: nil,
			expectedErr:   errExpected,
			errContains:   "rows error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.setupMock(mock)

			names, err := service.GetNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
				if tc.errContains != "" {
					require.ErrorContains(t, err, tc.errContains)
				}

				require.Nil(t, names)
			} else {
				if tc.errContains != "" {
					require.Error(t, err)
					require.ErrorContains(t, err, tc.errContains)
					require.Nil(t, names)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedNames, names)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			name: "success with unique names",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Ivan").
					AddRow("John").
					AddRow("Ivan")
				mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)

				return rows
			},
			expectedNames: []string{"Ivan", "John", "Ivan"},
			expectedErr:   nil,
		},
		{
			name: "success with empty result",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)
				return rows
			},
			expectedNames: []string{},
			expectedErr:   nil,
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				mock.ExpectQuery(queryUniqueNames).WillReturnError(errExpected)
				return nil
			},
			expectedNames: nil,
			expectedErr:   errExpected,
			errContains:   "db query",
		},
		{
			name: "scan error (nil value)",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)
				return rows
			},
			expectedNames: nil,
			expectedErr:   nil,
			errContains:   "rows scanning",
		},
		{
			name: "rows error",
			setupMock: func(mock sqlmock.Sqlmock) *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Bob").
					RowError(0, errExpected)
				mock.ExpectQuery(queryUniqueNames).WillReturnRows(rows)
				return rows
			},
			expectedNames: nil,
			expectedErr:   errExpected,
			errContains:   "rows error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			service := db.New(mockDB)

			tc.setupMock(mock)

			names, err := service.GetUniqueNames()

			if tc.expectedErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedErr)
				if tc.errContains != "" {
					require.ErrorContains(t, err, tc.errContains)
				}

				require.Nil(t, names)
			} else {
				if tc.errContains != "" {
					require.Error(t, err)
					require.ErrorContains(t, err, tc.errContains)
					require.Nil(t, names)
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedNames, names)
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}
