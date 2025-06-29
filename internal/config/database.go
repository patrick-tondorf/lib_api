package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewSupabaseDB() (*pgx.Conn, error) {

	config, _ := pgx.ParseConfig("")
	config.Host = os.Getenv("HOST")
	config.User = os.Getenv("USER")
	config.Password = os.Getenv("PASSWORD")
	config.Database = os.Getenv("DBNAME")
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetSecretKey() string {
	return os.Getenv("SECRET_KEY")
}
