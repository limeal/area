package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db *redis.Client
}

// RedisConfig is a struct that contains the configuration for a Redis connection.
// @property {string} Host - The hostname of the Redis server.
// @property {int} Port - The port to connect to Redis on.
// @property {string} Username - The username to use when connecting to the Redis server.
// @property {string} Password - The password to use when connecting to the Redis server.
// @property {int} Database - The database number to connect to.
// @property {string} URL - The URL of the Redis server.
// @property {bool} Reset - If true, the database will be flushed before the application starts.
// @property TLSConfig - This is a pointer to a tls.Config struct. This is used to configure the TLS
// connection to the Redis server.
// @property {int} PoolSize - The size of the connection pool.
type RedisConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Database  int
	URL       string
	Reset     bool
	TLSConfig *tls.Config
	PoolSize  int
}

// New creates a new redis storage
func CreateRedisStorage() fiber.Storage {

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:       0,
		Username: "",
		Password: "",
	})

	// Flush all keys (Comment this line if you don't want to flush all keys)
	db.FlushAll(context.Background())

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Successful connection to redis")
	return &Storage{db: db}
}

// Getting the value of the key from the redis database.
func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val, err := s.db.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Setting the key and value in the redis database.
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return s.db.Set(context.Background(), key, val, exp).Err()
}

// Deleting the key from the redis database.
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	return s.db.Del(context.Background(), key).Err()
}

// Flushing the database.
func (s *Storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

// Closing the connection to the redis database.
func (s *Storage) Close() error {
	return s.db.Close()
}

// This is a function that returns the redis client.
func (s *Storage) Conn() *redis.Client {
	return s.db
}
