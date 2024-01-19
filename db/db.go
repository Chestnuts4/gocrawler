package db

import (
	"database/sql"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
)

// DB is a wrapper around a SQL database and a Redis database.
type DB struct {
	SQLDB   *sql.DB
	RedisDB *redis.Client
}
