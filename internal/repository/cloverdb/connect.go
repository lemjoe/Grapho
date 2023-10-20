package cloverdb

import "github.com/ostafen/clover/v2"

type DB struct {
	db *clover.DB
}

func ConnectDB(dir string) (*DB, error) {
	db, err := clover.Open("db")
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}
func (d *DB) Close() {
	d.db.Close()
}
