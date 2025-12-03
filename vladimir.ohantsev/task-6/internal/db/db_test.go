package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	ErrSome = errors.New("some error")
)

type testList struct {
	name          string
	values        []string
	expectedError error
}

func ListMock(name string, values []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{name})
	for _, name := range values {
		rows = rows.AddRow(name)
	}

	return rows
}

func testGetNames(
	t *testing.T,
	testFunc func(service DBService) ([]string, error),
	query string,
) {
	tests := []testList{
		{
			name:   "success case",
			values: []string{"Ivan", "Gena228"},
		},
		{
			name:   "empty case",
			values: nil,
		},
		{
			name:          "error case",
			values:        nil,
			expectedError: ErrSome,
		},
	}

	db, mock, err := sqlmock.New()
	require.NoError(
		t, err,
		"failed to create sqlmock: %s", err,
	)

	defer db.Close()

	service := DBService{DB: db}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock.ExpectQuery(query).
				WillReturnRows(ListMock("name", test.values)).
				WillReturnError(test.expectedError)

			names, err := testFunc(service)

			require.Equal(
				t, test.values, names,
				"expected: %s, actual: %s", test.values, names,
			)

			if test.expectedError != nil {
				require.ErrorIs(
					t, err, test.expectedError,
					"expected: %s, actual: %s", test.expectedError, err,
				)
			} else {
				require.NoError(
					t, err,
					"unexpected error: %s", err,
				)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
		})
	}

	t.Run("rows error", func(t *testing.T) {
		mock.ExpectQuery(query).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"name"}).
					AddRow(nil).
					RowError(0, ErrSome),
			)

		names, err := testFunc(service)

		require.ErrorIs(
			t, err, ErrSome,
			"expected: %s, actual: %s", ErrSome, err,
		)

		require.Nil(
			t, names,
			"names not nil: %s", names,
		)

		err = mock.ExpectationsWereMet()
		require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
	})

	t.Run("scan error", func(t *testing.T) {
		mock.ExpectQuery(query).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"name"}).
					AddRow(nil),
			)

		names, err := testFunc(service)

		require.Error(
			t, err,
			"expected error",
		)

		require.Nil(
			t, names,
			"names not nil: %s", names,
		)

		err = mock.ExpectationsWereMet()
		require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
	})
}

func TestGetNames(t *testing.T) {
	testGetNames(t, DBService.GetNames, "SELECT name FROM users")
}

func TestGetUniqueNames(t *testing.T) {
	testGetNames(t, DBService.GetUniqueNames, "SELECT DISTINCT name FROM users")
}

func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(
		t, err,
		"failed to create sqlmock: %s", err,
	)

	defer db.Close()

	service := New(db)
	require.Equal(t, service.DB, db, "use another object in service")
}
