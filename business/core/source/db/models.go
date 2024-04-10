package db

type Campaign struct {
	ID     string `db:"id"`
	Name   string `db:"name"`
	Domain string `db:"domain"`
}

type Source struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
