package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

type Mysql struct {
	db *sql.DB
}

type URLFilter struct {
	URL   string
	Alias string
	ID    int64
}

func (db *Mysql) CreateURL(sURL string, alias string) (int64, error) {
	smt, err := db.db.Prepare("INSERT INTO url (url, Alias) VALUES (?, ?)")
	if err != nil {
		return 0, errors.Wrap(err, "url prepare error")
	}
	res, err := smt.Exec(sURL, alias)
	if err != nil {
		var dbErr *mysql.MySQLError
		if errors.As(err, &dbErr) {
			if dbErr.Number == 1062 {
				return 0, errors.New("url already exists")
			}
		}
		return 0, errors.Wrap(err, "create url error")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(err, "last insert id error")
	}
	return id, nil
}

type URL struct {
	id    int64
	url   string
	alias string
}

func (db *Mysql) GetURL(filter URLFilter) (*URL, error) {
	sqlQuery := "SELECT * FROM url WHERE 1=1 "
	if filter.URL != "" {
		sqlQuery += "AND url = '" + filter.URL + "' "
	}

	if filter.Alias != "" {
		sqlQuery += "AND Alias = '" + filter.Alias + "' "
	}

	if filter.ID > 0 {
		sqlQuery += "AND id = " + fmt.Sprint(filter.ID) + " "
	}

	smt, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return nil, errors.Wrap(err, "get url prepare error")
	}
	res, err := smt.Query()
	if err != nil {
		return nil, errors.Wrap(err, "get url error")
	}

	if res.Next() {
		var id int64
		var url string
		var alias string
		err = res.Scan(&id, &url, &alias)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}
		return &URL{id: id, url: url, alias: alias}, nil
	}
	return nil, errors.New("url not found")
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
