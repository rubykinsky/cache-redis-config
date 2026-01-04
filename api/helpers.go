package helpers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

const (
	defaultAddr     = "localhost:6379"
	defaultPassword = ""
	defaultDB       = 0
)

func getRedisConfig() (string, string, int) {
	addr := viper.GetString("redis.addr")
	if addr == "" {
		addr = defaultAddr
	}
	password := viper.GetString("redis.password")
	if password == "" {
		password = defaultPassword
	}
	db, err := strconv.Atoi(viper.GetString("redis.db"))
	if err != nil || db < 0 {
		db = defaultDB
	}
	return addr, password, db
}

func newRedisClient(ctx context.Context, addr string, password string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := redisClient.Ping(ctx).Result()
	return redisClient, err
}

func newRedisClientWithTLS(ctx context.Context, addr string, password string, db int, tlsConfig *tls.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTLS:      tlsConfig,
		TLSConfig:   tlsConfig,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout:  30 * time.Second,
		MinIdleConns:  10,
		MaxConnAge:    5 * time.Minute,
		MaxRetries:    3,
	})
	_, err := redisClient.Ping(ctx).Result()
	return redisClient, err
}

func loadTLSCertificates(caCertPath string, clientCertPath string, clientKeyPath string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
		MinVersion:  tls.VersionTLS12,
	}
	return tlsConfig, nil
}

func getRedisClient(ctx context.Context) (*redis.Client, error) {
	addr, password, db := getRedisConfig()
	tlsEnabled := viper.GetBool("redis.tls")
	if tlsEnabled {
		caCertPath := viper.GetString("redis.tls.caCertPath")
		clientCertPath := viper.GetString("redis.tls.clientCertPath")
		clientKeyPath := viper.GetString("redis.tls.clientKeyPath")
		tlsConfig, err := loadTLSCertificates(caCertPath, clientCertPath, clientKeyPath)
		if err != nil {
			return nil, err
		}
		return newRedisClientWithTLS(ctx, addr, password, db, tlsConfig)
	}
	return newRedisClient(ctx, addr, password, db)
}