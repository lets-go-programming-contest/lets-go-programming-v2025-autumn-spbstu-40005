package db_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/P3rCh1/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

const (
	queryGetNames       = "SELECT name FROM users"
	queryGetUniqueNames = "SELECT DISTINCT name FROM users"
)

var ErrSome = errors.New("some error")

type testcase struct {
	name          string
	values        []string
	expectedError error
}

var tests = []testcase{ //nolint:gochecknoglobals
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

func helperListMock(t *testing.T, name string, values []string) *sqlmock.Rows {
	t.Helper()

	rows := sqlmock.NewRows([]string{name})
	for _, name := range values {
		rows = rows.AddRow(name)
	}

	return rows
}

func helperInitMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) { //nolint:ireturn
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(
		t, err,
		"failed to init sqlmock: %s", err,
	)

	return db, mock
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	for _, test := range tests { //nolint:paralleltest
		t.Run(test.name, func(t *testing.T) {
			testGetNamesPredicted(t, &test)
		})
	}

	t.Run("special case: scan error", testGetNamesScanError) //nolint:paralleltest
	t.Run("special case: row error", testGetNamesRowsError)  //nolint:paralleltest
}

func testGetNamesPredicted(t *testing.T, test *testcase) { //nolint:thelper
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(helperListMock(t, "name", test.values)).
		WillReturnError(test.expectedError)

	names, err := service.GetNames()

	require.Equal(t, test.values, names)

	if test.expectedError != nil {
		require.ErrorIs(t, err, test.expectedError)
	} else {
		require.NoError(t, err)
	}

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}

func testGetNamesRowsError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow("egor").
				RowError(0, ErrSome),
		)

	names, err := service.GetNames()

	require.ErrorIs(t, err, ErrSome)
	require.Nil(t, names)

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}

func testGetNamesScanError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow(nil),
		)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	for _, test := range tests { //nolint:paralleltest
		t.Run(test.name, func(t *testing.T) {
			testGetUniqueNamesPredicted(t, &test)
		})
	}

	t.Run("special case: scan error", testGetUniqueNamesScanError) //nolint:paralleltest
	t.Run("special case: row error", testGetUniqueNamesRowsError)  //nolint:paralleltest
}

func testGetUniqueNamesPredicted(t *testing.T, test *testcase) { //nolint:thelper
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(helperListMock(t, "name", test.values)).
		WillReturnError(test.expectedError)

	names, err := service.GetUniqueNames()

	require.Equal(t, test.values, names)

	if test.expectedError != nil {
		require.ErrorIs(t, err, test.expectedError)
	} else {
		require.NoError(t, err)
	}

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}

func testGetUniqueNamesRowsError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow("egor").
				RowError(0, ErrSome),
		)

	names, err := service.GetUniqueNames()

	require.ErrorIs(t, err, ErrSome)
	require.Nil(t, names)

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}

func testGetUniqueNamesScanError(t *testing.T) {
	t.Parallel()

	db, mock := helperInitMock(t)
	defer db.Close()

	service := database.New(db)

	mock.ExpectQuery(queryGetUniqueNames).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"name"}).
				AddRow(nil),
		)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)

	err = mock.ExpectationsWereMet()
	require.NoError(t, mock.ExpectationsWereMet(), "expectations won't met: %s", err)
}
