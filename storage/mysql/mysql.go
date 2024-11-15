package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	"log"
)

type Mysql struct {
	db *sql.DB
}

func Connection(storagePath string) (*Mysql, error) {
	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		return nil, errors.Wrap(err, "db open error")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping error")
	}

	return &Mysql{db: db}, nil
}

func ConnectionClose(db *Mysql) {
	err := db.db.Close()
	if err != nil {
		log.Fatal(errors.Wrap(err, "db close error"))
	}
}
