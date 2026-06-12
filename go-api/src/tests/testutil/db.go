package testutil

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/jordanhuaman/go-api/src/models"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func detectDockerHost() string {
	if h := os.Getenv("DOCKER_HOST"); h != "" {
		return h
	}
	out, err := exec.Command("docker", "context", "inspect", "--format", "{{.Endpoints.docker.Host}}").Output()
	if err == nil {
		host := strings.TrimSpace(string(out))
		if host != "" {
			return host
		}
	}
	return ""
}

func openGORM(t *testing.T, dsn string) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Statistics{},
		&models.MatrixInput{},
		&models.MatrixResult{},
		&models.RefreshToken{},
	); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	return db
}

func resolveDockerPort(t *testing.T, c *tcpostgres.PostgresContainer) string {
	t.Helper()

	id := c.GetContainerID()
	if id == "" {
		t.Fatal("failed to get container ID")
	}

	out, err := exec.Command("docker", "port", id, "5432/tcp").Output()
	if err != nil {
		t.Fatalf("failed to get port via docker CLI: %v", err)
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ":")
	return strings.TrimSpace(parts[len(parts)-1])
}

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	if dsn := os.Getenv("TEST_DATABASE_URL"); dsn != "" {
		t.Log(fmt.Sprintf("using TEST_DATABASE_URL: %s", dsn))
		return openGORM(t, dsn)
	}

	if err := exec.Command("docker", "info").Run(); err != nil {
		t.Skip("Docker is not available, skipping integration test")
	}

	if host := detectDockerHost(); host != "" {
		t.Log(fmt.Sprintf("detected DOCKER_HOST: %s", host))
		t.Setenv("DOCKER_HOST", host)
	}

	ctx := context.Background()

	pgContainer, err := tcpostgres.RunContainer(ctx,
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate postgres container: %v", err)
		}
	})

	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	var portStr string
	mappedPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Log("MappedPort failed, falling back to docker CLI")
		portStr = resolveDockerPort(t, pgContainer)
	} else {
		portStr = mappedPort.Port()
	}

	connStr := fmt.Sprintf("postgres://test:test@%s/%s?search_path=public&sslmode=disable", net.JoinHostPort(host, portStr), "testdb")
	t.Log(fmt.Sprintf("test database ready at %s:%s", host, portStr))
	return openGORM(t, connStr)
}
