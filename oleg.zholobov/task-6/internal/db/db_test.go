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

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	scenarios := []testScenario{
		{
			description:  "regular query execution",
			expectedData: []string{"Dmitry", "Anna", "Sergey"},
		},
		{
			description: "query execution failure",
			dbError:     errDatabaseFailure,
			shouldFail:  true,
			failReason:  "query",
		},
		{
			description:  "single record result",
			expectedData: []string{"Ekaterina"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

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
					assert.Contains(t, err.Error(), "db query")
				}

				assert.Nil(t, results, "results should be nil on error")
			} else {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, scenario.expectedData, results)
			}

			assert.NoError(t, sqlMock.ExpectationsWereMet(),
				"unfulfilled sqlmock expectations")
		})
	}

	t.Run("scan error - null value handling", func(t *testing.T) {
		t.Parallel()

		sqlDB, sqlMock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		sqlMock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		results, err := service.GetNames()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.Nil(t, results)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("row iteration error", func(t *testing.T) {
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
		assert.Contains(t, err.Error(), "rows error")
		assert.Nil(t, results)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("empty result set - nil slice", func(t *testing.T) {
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
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	scenarios := []testScenario{
		{
			description:  "distinct query execution",
			expectedData: []string{"Dmitry", "Anna", "Sergey"},
		},
		{
			description: "distinct query failure",
			dbError:     errDatabaseFailure,
			shouldFail:  true,
			failReason:  "query",
		},
		{
			description:  "single distinct record",
			expectedData: []string{"Ekaterina"},
		},
		{
			description:  "duplicate names in distinct query",
			expectedData: []string{"Dmitry", "Dmitry", "Anna"},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			t.Parallel()

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
					assert.Contains(t, err.Error(), "db query")
				}

				assert.Nil(t, results, "results should be nil on error")
			} else {
				require.NoError(t, err, "unexpected error")
				assert.Equal(t, scenario.expectedData, results)
			}

			assert.NoError(t, sqlMock.ExpectationsWereMet(),
				"unfulfilled sqlmock expectations")
		})
	}

	t.Run("distinct scan error - null handling", func(t *testing.T) {
		t.Parallel()

		sqlDB, sqlMock, err := sqlmock.New()
		require.NoError(t, err)
		defer sqlDB.Close()

		service := db.New(sqlDB)

		sqlMock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		results, err := service.GetUniqueNames()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.Nil(t, results)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("distinct row iteration error", func(t *testing.T) {
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
		assert.Contains(t, err.Error(), "rows error")
		assert.Nil(t, results)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("distinct empty result - nil slice", func(t *testing.T) {
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
	})
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
