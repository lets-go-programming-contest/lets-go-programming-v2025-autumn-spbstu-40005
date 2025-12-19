package db_test

import (
	"errors"
	"testing"

	"anton.mezentsev/task-6/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

const (
	querySelectAll      = "SELECT login FROM users"
	querySelectDistinct = "SELECT DISTINCT login FROM users"
)

var ErrTest = errors.New("test error")

func TestFetchAllLogins_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).
		AddRow("maxim").
		AddRow("sofia")
	mock.ExpectQuery(querySelectAll).WillReturnRows(rows)

	logins, err := service.FetchAllLogins()
	require.NoError(t, err)
	require.Equal(t, []string{"maxim", "sofia"}, logins)
}

func TestFetchAllLogins_QueryFails(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	mock.ExpectQuery(querySelectAll).WillReturnError(ErrTest)

	logins, err := service.FetchAllLogins()
	require.ErrorContains(t, err, "failed to execute query")
	require.Nil(t, logins)
}

func TestFetchAllLogins_BadRowData(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).AddRow(nil)
	mock.ExpectQuery(querySelectAll).WillReturnRows(rows)

	logins, err := service.FetchAllLogins()
	require.ErrorContains(t, err, "failed to scan row")
	require.Nil(t, logins)
}

func TestFetchAllLogins_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).AddRow("alex")
	rows.RowError(0, ErrTest)
	mock.ExpectQuery(querySelectAll).WillReturnRows(rows)

	logins, err := service.FetchAllLogins()
	require.ErrorContains(t, err, "row iteration error")
	require.Nil(t, logins)
}

func TestFetchUniqueLogins_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).
		AddRow("nikita").
		AddRow("elena")
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	logins, err := service.FetchUniqueLogins()
	require.NoError(t, err)
	require.Equal(t, []string{"nikita", "elena"}, logins)
}

func TestFetchUniqueLogins_QueryFails(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	mock.ExpectQuery(querySelectDistinct).WillReturnError(ErrTest)

	logins, err := service.FetchUniqueLogins()
	require.ErrorContains(t, err, "failed to execute query")
	require.Nil(t, logins)
}

func TestFetchUniqueLogins_ScanFails(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).AddRow(nil)
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	logins, err := service.FetchUniqueLogins()
	require.ErrorContains(t, err, "failed to scan row")
	require.Nil(t, logins)
}

func TestFetchUniqueLogins_RowError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.UserRepository{DB: mockDB}

	rows := sqlmock.NewRows([]string{"login"}).AddRow("dmitry")
	rows.RowError(0, ErrTest)
	mock.ExpectQuery(querySelectDistinct).WillReturnRows(rows)

	logins, err := service.FetchUniqueLogins()
	require.ErrorContains(t, err, "row iteration error")
	require.Nil(t, logins)
}

func TestCreateRepository(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	repo := db.CreateRepository(mockDB)
	require.Equal(t, mockDB, repo.DB)
}
