package database

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	re := regexp.MustCompile(`^(.*` + "orders" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		return err
	}

	return nil
}

func Pool() (*pgxpool.Pool, error) {
	err := LoadEnv()

	if err != nil {
		return nil, err
	}

	uri := os.Getenv("DB_URI")

	if uri == "" {
		return nil, fmt.Errorf("uri not found")
	}

	config, err := pgxpool.ParseConfig(uri)

	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
