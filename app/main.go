package main

import (
	_ "embed"
	"fmt"
	"github.com/grumpycatyo-collab/turbo-pancake/business/data/dbschema"
	"github.com/grumpycatyo-collab/turbo-pancake/business/sys/database"
)

type DB struct {
	User       string `conf:"default:admin"`
	Password   string `conf:"default:admin,mask"`
	Host       string `conf:"default:localhost"`
	Name       string `conf:"default:db"`
	DisableTLS bool   `conf:"default:true"`
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("\nStartup error \n")
		panic("Server crushed")
	}
}

func run() error {

	config := DB{
		User:       "admin",
		Password:   "admin",
		Host:       "localhost",
		Name:       "db",
		DisableTLS: true,
	}

	db, err := database.Open(database.Config{
		User:     config.User,
		Password: config.Password,
		Host:     config.Host,
		Name:     config.Name,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer db.Close()

	if err := dbschema.DropAll(db); err != nil {
		fmt.Printf("create tables: %w", err)
	}

	if err := dbschema.Create(db); err != nil {
		fmt.Printf("create tables: %w", err)
	}

	if err := dbschema.Seed(db); err != nil {
		fmt.Printf("seed database: %w", err)
	}

	fmt.Printf("Data seed successful \n")

	return nil
}
