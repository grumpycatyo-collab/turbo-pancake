package db

type Campaign struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type Source struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
