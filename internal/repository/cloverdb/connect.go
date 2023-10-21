package cloverdb

import "github.com/ostafen/clover/v2"

type DB struct {
	DB *clover.DB
}

func ConnectDB(dir string) (*DB, error) {
	db, err := clover.Open(dir)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB: db,
	}, nil
}
func (d *DB) Close() {
	d.DB.Close()
}
