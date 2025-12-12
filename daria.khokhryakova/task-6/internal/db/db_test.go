package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DariaKhokhryakova/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	selectNamesQuery       = "SELECT name FROM users"
	selectUniqueNamesQuery = "SELECT DISTINCT name FROM users"
)

type mockDatabase struct {
	mock.Mock
}

func (m *mockDatabase) Query(query string, args ...any) (*sql.Rows, error) {
	argsList := m.Called(query, args)
	return argsList.Get(0).(*sql.Rows), argsList.Error(1)
}

func TestDBService_GetNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).
		AddRow("John").
		AddRow("Gerald").
		AddRow("Michael")

	sqlMock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Gerald", "Michael"}, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetNames_DBError(t *testing.T) {
	t.Parallel()

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectNamesQuery, []any{}).Return((*sql.Rows)(nil), assert.AnError)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetNames_ScanFailure(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).AddRow(123)

	sqlMock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetNames_RowIterationError(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).
		AddRow("John").
		RowError(0, assert.AnError)

	sqlMock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetNames_NoData(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns)

	sqlMock.ExpectQuery(selectNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{}, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).
		AddRow("John").
		AddRow("Gerald").
		AddRow("John").
		AddRow("Michael").
		AddRow("Gerald")

	sqlMock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectUniqueNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Gerald", "John", "Michael", "Gerald"}, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetUniqueNames_DBError(t *testing.T) {
	t.Parallel()

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectUniqueNamesQuery, []any{}).Return((*sql.Rows)(nil), assert.AnError)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetUniqueNames_ScanFailure(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).AddRow(123)

	sqlMock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectUniqueNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetUniqueNames_RowIterationError(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns).
		AddRow("John").
		RowError(0, assert.AnError)

	sqlMock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectUniqueNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_GetUniqueNames_NoData(t *testing.T) {
	t.Parallel()

	dbConn, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbConn.Close()

	columns := []string{"name"}
	rows := sqlmock.NewRows(columns)

	sqlMock.ExpectQuery(selectUniqueNamesQuery).WillReturnRows(rows)

	mockDB := &mockDatabase{}
	mockDB.On("Query", selectUniqueNamesQuery, []any{}).Return(rows, nil)

	service := db.New(mockDB)
	result, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, []string{}, result)
	mockDB.AssertExpectations(t)
}

func TestDBService_Initialization(t *testing.T) {
	t.Parallel()

	mockDB := &mockDatabase{}
	service := db.New(mockDB)

	assert.NotNil(t, service)
	assert.Equal(t, mockDB, service.DB)
}
