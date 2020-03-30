package db

import (
	"database/sql"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/johnshiver/asapp_challenge/config"
)

var (
	singletonDB *sqlx.DB
	dbOnce      sync.Once
)

func Get() (*sqlx.DB, error) {
	dbOnce.Do(func() {
		sqlDB, err := sql.Open("postgres", config.Get().DBConn)
		if err != nil {
			log.Fatalf("couldn't connect to db: %v", err)
		}
		singletonDB = sqlx.NewDb(sqlDB, "postgres")
	})
	if err := singletonDB.Ping(); err != nil {
		log.Printf("couldn't ping db: %v", err)
		return nil, err
	}
	return singletonDB, nil
}
