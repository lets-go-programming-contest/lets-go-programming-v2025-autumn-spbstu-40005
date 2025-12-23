package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "oleg.zholobov/task-6/internal/db"
)

var (
	errDatabaseFailure = errors.New("database connection failure")
	errRowProcessing   = errors.New("row processing error")
)

type testScenario struct {
	description  string
	dbError      error
	expectedData []string
	shouldFail   bool
	failReason   string // "query", "scan", "rows"
}

func TestDBService_GetNames_RegularQueryExecution(t *testing.T) {
	t.Parallel()

	testDBService_GetNames(t, testScenario{
		description:  "regular query execution",
		expectedData: []string{"Dmitry", "Anna", "Sergey"},
	})
}

func TestDBService_GetNames_QueryExecutionFailure(t *testing.T) {
	t.Parallel()

	testDBService_GetNames(t, testScenario{
		description: "query execution failure",
		dbError:     errDatabaseFailure,
		shouldFail:  true,
		failReason:  "query",
	})
}

func TestDBService_GetNames_SingleRecordResult(t *testing.T) {
	t.Parallel()

	testDBService_GetNames(t, testScenario{
		description:  "single record result",
		expectedData: []string{"Ekaterina"},
	})
}

func testDBService_GetNames(t *testing.T, scenario testScenario) {
	t.Helper()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err, "sqlmock initialization failed")
	defer sqlDB.Close()

	service := db.New(sqlDB)

	if scenario.dbError != nil {
		sqlMock.ExpectQuery("SELECT name FROM users").
			WillReturnError(scenario.dbError)
	} else {
		resultRows := sqlmock.NewRows([]string{"name"})
		for _, name := range scenario.expectedData {
			resultRows.AddRow(name)
		}

		sqlMock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(resultRows)
	}

	results, err := service.GetNames()

	if scenario.shouldFail {
		require.Error(t, err, "expected error but got none")

		if scenario.failReason == "query" {
			assert.ErrorContains(t, err, "db query")
		}

		assert.Nil(t, results, "results should be nil on error")
	} else {
		require.NoError(t, err, "unexpected error")
		assert.Equal(t, scenario.expectedData, results)
	}

	assert.NoError(t, sqlMock.ExpectationsWereMet(),
		"failed sqlmock expectations")
}

func TestDBService_GetNames_ScanError_NullValueHandling(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	sqlMock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	results, err := service.GetNames()
	require.Error(t, err)
	assert.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowIterationError(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	testRows := sqlmock.NewRows([]string{"name"}).AddRow("Vladimir")
	testRows.RowError(0, errRowProcessing)
	sqlMock.ExpectQuery("SELECT name FROM users").WillReturnRows(testRows)

	results, err := service.GetNames()
	require.Error(t, err)
	assert.ErrorContains(t, err, "rows error")
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_GetNames_EmptyResultSet_NilSlice(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	sqlMock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}))

	results, err := service.GetNames()
	require.NoError(t, err)
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_DistinctQueryExecution(t *testing.T) {
	t.Parallel()

	testDBService_GetUniqueNames(t, testScenario{
		description:  "distinct query execution",
		expectedData: []string{"Dmitry", "Anna", "Sergey"},
	})
}

func TestDBService_GetUniqueNames_DistinctQueryFailure(t *testing.T) {
	t.Parallel()

	testDBService_GetUniqueNames(t, testScenario{
		description: "distinct query failure",
		dbError:     errDatabaseFailure,
		shouldFail:  true,
		failReason:  "query",
	})
}

func TestDBService_GetUniqueNames_SingleDistinctRecord(t *testing.T) {
	t.Parallel()

	testDBService_GetUniqueNames(t, testScenario{
		description:  "single distinct record",
		expectedData: []string{"Ekaterina"},
	})
}

func TestDBService_GetUniqueNames_DuplicateNamesInDistinctQuery(t *testing.T) {
	t.Parallel()

	testDBService_GetUniqueNames(t, testScenario{
		description:  "duplicate names in distinct query",
		expectedData: []string{"Dmitry", "Dmitry", "Anna"},
	})
}

func testDBService_GetUniqueNames(t *testing.T, scenario testScenario) {
	t.Helper()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err, "sqlmock initialization failed")
	defer sqlDB.Close()

	service := db.New(sqlDB)

	if scenario.dbError != nil {
		sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(scenario.dbError)
	} else {
		resultRows := sqlmock.NewRows([]string{"name"})
		for _, name := range scenario.expectedData {
			resultRows.AddRow(name)
		}

		sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(resultRows)
	}

	results, err := service.GetUniqueNames()

	if scenario.shouldFail {
		require.Error(t, err, "expected error but got none")

		if scenario.failReason == "query" {
			assert.ErrorContains(t, err, "db query")
		}

		assert.Nil(t, results, "results should be nil on error")
	} else {
		require.NoError(t, err, "unexpected error")
		assert.Equal(t, scenario.expectedData, results)
	}

	assert.NoError(t, sqlMock.ExpectationsWereMet(),
		"unmet sqlmock expectations")
}

func TestDBService_GetUniqueNames_DistinctScanError_NullHandling(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

	results, err := service.GetUniqueNames()
	require.Error(t, err)
	assert.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_DistinctRowIterationError(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	testRows := sqlmock.NewRows([]string{"name"}).AddRow("Vladimir")
	testRows.RowError(0, errRowProcessing)
	sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(testRows)

	results, err := service.GetUniqueNames()
	require.Error(t, err)
	assert.ErrorContains(t, err, "rows error")
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_DistinctEmptyResult_NilSlice(t *testing.T) {
	t.Parallel()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	service := db.New(sqlDB)

	sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"name"}))

	results, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Nil(t, results)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestDBService_Constructor(t *testing.T) {
	t.Parallel()

	sqlDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer sqlDB.Close()

	dbService := db.New(sqlDB)
	require.NotNil(t, dbService)
	assert.Equal(t, sqlDB, dbService.DB, "database connection mismatch")
}
