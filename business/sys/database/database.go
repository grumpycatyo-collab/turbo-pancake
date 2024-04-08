package database

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"reflect"
	"time"
)

type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

var (
	ErrDBNotFound = errors.New("not found")
)

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

func NamedQuerySlice(log *zerolog.Logger, db *sqlx.DB, query string, data interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	fmt.Printf("%v", data)
	rows, err := sqlx.NamedQuery(db, query, data)
	if err != nil {
		return err
	}

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}
