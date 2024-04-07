package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

func Open(cfg Config) (*sqlx.DB, error) {
	dsnString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", cfg.User, cfg.Password, cfg.Host, "3306", cfg.Name)

	db, err := sqlx.Open("mysql", dsnString)
	if err != nil {
		fmt.Printf("connect to db: %s", err)
		return nil, err
	}
	return db, nil
}

func StatusCheck(db *sqlx.DB) error {
	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	const q = `SELECT true`
	var tmp bool
	return db.QueryRow(q).Scan(&tmp)
}
