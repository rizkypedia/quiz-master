package tests

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github/rizkypedia/quiz-master/internal/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mockDatabase(t *testing.T) sqlmock.Sqlmock {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		_ = sqlDB.Close()
		t.Fatalf("open GORM with SQL mock: %v", err)
	}

	originalDB := database.DB
	database.DB = gormDB
	t.Cleanup(func() {
		database.DB = originalDB
		_ = sqlDB.Close()
	})

	return mock
}
