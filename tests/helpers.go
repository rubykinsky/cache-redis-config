package helpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis/v9"
)

type config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func loadConfig(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c config
	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func newRedisClient(c *config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Host + ":" + c.Port,
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func getCacheDir() string {
	cacheDir := os.Getenv("CACHE_DIR")
	if cacheDir == "" {
		cacheDir = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return cacheDir
}

func getRedisConfigPath() string {
	redisConfigPath := os.Getenv("REDIS_CONFIG_PATH")
	if redisConfigPath == "" {
		redisConfigPath = filepath.Join(getCacheDir(), "redis.json")
	}
	return redisConfigPath
}

var ctx = context.Background()

func init() {
	if _, err := os.Stat(getRedisConfigPath()); os.IsNotExist(err) {
		log.Fatal("Redis config file not found")
	}
}