package repository

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

const (
	testLocalDockerUbuntu  = "unix:///home/olegvantsevich/.docker/desktop/docker.sock"
	testLocalDockerWindows = ""

	testMigrationPassUbuntu  = "/home/olegvantsevich/GolandProjects/Payment-Service"
	testMigrationPassWindows = "C:/Users/oleg/GolandProjects/Payment-Service"

	testPostgresPort     = "4444"
	testPostgresHost     = "localhost"
	testPostgresName     = "postgres-test"
	testPostgresDB       = "postgres"
	testPostgresUser     = "postgres"
	testPostgresPassword = "postgres"
)

var testAccountRepository *Account
var testTransactionRepository *Transaction
var testTransactor PgxTransactor

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool(testLocalDockerWindows)
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	postgres, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       testPostgresName,
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", testPostgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPostgresPassword),
			fmt.Sprintf("POSTGRES_DB=%s", testPostgresDB),
			"listen_addresses = '*'",
		},
		Mounts: []string{fmt.Sprintf("%s/migrations:/docker-entrypoint-initdb.d", testMigrationPassWindows)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: testPostgresHost, HostPort: fmt.Sprintf("%s/tcp", testPostgresPort)}},
		},
	},
		func(config *docker.HostConfig) {
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
	postgres.Expire(60)

	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	ctx := context.Background()
	if err = pool.Retry(func() error {
		pgPool, retryErr := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", testPostgresUser, testPostgresPassword, testPostgresHost, testPostgresPort, testPostgresDB))
		if retryErr != nil {
			return fmt.Errorf("could not connect to db %s", retryErr)
		}
		retryErr = pgPool.Ping(ctx)
		if retryErr != nil {
			return retryErr
		}
		testTransactionRepository = NewTransaction(NewPgxWithinTransactionRunner(pgPool))
		testAccountRepository = NewAccount(NewPgxWithinTransactionRunner(pgPool))
		testTransactor = NewPgxTransactor(pgPool)
		return nil
	}); err != nil {
		logrus.Fatalf("Could not connect to postgres: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(postgres); err != nil {
		logrus.Fatalf("Could not purge postgres: %s", err)
	}

	os.Exit(code)
}
