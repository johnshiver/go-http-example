package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type Config struct {
	DBConn        string        `default:"postgres://asapp:asapp@db:5432/asapp_chat?sslmode=disable"`
	TestDBConn    string        `default:"postgres://test:test@localhost:3345/asapp_chat_test?sslmode=disable"`
	ServerPort    string        `default:":8080"`
	JwtSecretKey  string        `default:"cant-touch-this"`
	JwtExpiration time.Duration `default:"10m"`
}

var (
	c          Config
	configOnce sync.Once
)

func Get() Config {
	configOnce.Do(func() {
		envconfig.MustProcess("ASAPP", &c)
	})
	return c
}
