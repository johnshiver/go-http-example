package testing

import (
	"database/sql"
	"log"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/johnshiver/asapp_challenge/config"
)

var (
	singletonDB *sqlx.DB
	dbOnce      sync.Once
)

func GetTestDB(t *testing.T) *sqlx.DB {
	dbOnce.Do(func() {
		sqlDB, err := sql.Open("postgres", config.Get().TestDBConn)
		if err != nil {
			log.Fatalf("couldn't connect to db: %v", err)
		}
		singletonDB = sqlx.NewDb(sqlDB, "postgres")
	})
	if err := singletonDB.Ping(); err != nil {
		log.Fatalf("couldn't ping db: %v", err)
	}
	return singletonDB
}
