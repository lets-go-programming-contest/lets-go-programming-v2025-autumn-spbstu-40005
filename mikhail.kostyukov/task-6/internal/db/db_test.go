package db_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	myDB "github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-6/internal/db"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)

	defer mockDB.Close()

	service := myDB.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}
