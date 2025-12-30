package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	myDB "github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var ErrExpected = errors.New("expected error")

type testCase struct {
	name          string
	mockSetup     func(sqlmock.Sqlmock)
	expectedNames []string
	expectError   bool
	errorContains string
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	defer func() { require.NoError(t, mockDB.Close()) }()

	service := myDB.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	runCommonTests(t, "SELECT name FROM users", func(s myDB.DBService) ([]string, error) {
		return s.GetNames()
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	runCommonTests(t, "SELECT DISTINCT name FROM users", func(s myDB.DBService) ([]string, error) {
		return s.GetUniqueNames()
	})
}

func getTestCases(query string) []testCase {
	return []testCase{
		{
			name: "success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Misha").AddRow("Masha")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedNames: []string{"Misha", "Masha"},
			expectError:   false,
			errorContains: "",
		},
		{
			name: "query error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(ErrExpected)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "db query",
		},
		{
			name: "scan error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows scanning",
		},
		{
			name: "rows iteration error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Misha").RowError(0, ErrExpected)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedNames: nil,
			expectError:   true,
			errorContains: "rows error",
		},
	}
}

func runCommonTests(t *testing.T, query string, callFunc func(myDB.DBService) ([]string, error)) {
	t.Helper()

	tests := getTestCases(query)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			defer func() { require.NoError(t, db.Close()) }()

			service := myDB.New(db)

			testCase.mockSetup(mock)

			names, err := callFunc(service)

			if testCase.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), testCase.errorContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
