package cloverdb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ostafen/clover/v2"
)

type DB struct {
	DB *clover.DB
}

func ConnectDB(dir string) (*DB, error) {
	//relative path to absolute
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("db path: %s\n", absDir)
	//check if dir exists
	_, err = os.Stat(absDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(absDir, 0740)
			if err != nil {
				//log.Print("Unable to create database: ", err)
				return nil, fmt.Errorf("unable to create database: %w", err)
			}
		} else {
			//log.Print("Unable to create database: ", err)
			return nil, fmt.Errorf("unable to create database: %w", err)
		}
	}

	db, err := clover.Open(absDir)
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
