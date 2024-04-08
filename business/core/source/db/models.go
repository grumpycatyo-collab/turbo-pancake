package db

// TODO: Optimize the using of two different model files
type Campaign struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type Source struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
