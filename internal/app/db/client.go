package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultHost     = "localhost"
	defaultPort     = 5432
	defaultUser     = "test"
	defaultPassword = "test"
	defaultDbName   = "test"
)

func NewTransactionManager(ctx context.Context) (*TransactionManager, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	return newTransactionManager(pool), nil
}

func getEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func generateDsn() string {
	host := getEnv("DB_HOST", defaultHost)
	port := getEnv("DB_PORT", strconv.Itoa(defaultPort))
	user := getEnv("DB_USER", defaultUser)
	password := getEnv("DB_PASSWORD", defaultPassword)
	dbname := getEnv("DB_DBNAME", defaultDbName)
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
