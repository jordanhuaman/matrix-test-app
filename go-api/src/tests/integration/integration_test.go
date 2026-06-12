package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/jordanhuaman/go-api/src/tests/testutil"
	"gorm.io/gorm"
)

var (
	jwtSec  = "test-secret-key-for-integration-tests"
	longTTL = 30 * 24 * time.Hour
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.SetupTestDB(t)
}
