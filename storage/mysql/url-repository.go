package mysql

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go_first/internal/lib/common/model/urlshort"
)

type URLFilter struct {
	URL   string
	Alias string
	ID    int64
}

func (db *Mysql) CreateURL(sURL string, alias string) (*urlshort.URLShort, error) {
	smt, err := db.db.Prepare("INSERT INTO url (url, alias) VALUES (?, ?)")
	if err != nil {
		return nil, errors.Wrap(err, "url prepare error")
	}
	res, err := smt.Exec(sURL, alias)
	if err != nil {
		var dbErr *mysql.MySQLError
		if errors.As(err, &dbErr) {
			if dbErr.Number == 1062 {
				return nil, errors.New("url already exists")
			}
		}
		return nil, errors.Wrap(err, "create url error")
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "last insert id error")
	}
	return &urlshort.URLShort{
		ID:    id,
		URL:   sURL,
		Alias: alias,
	}, nil
}

func (db *Mysql) GetURL(filter URLFilter) (*urlshort.URLShort, error) {
	sqlQuery := ""
	if filter.URL != "" {
		sqlQuery += " AND url = '" + filter.URL + "' "
	}

	if filter.Alias != "" {
		sqlQuery += " AND alias = '" + filter.Alias + "' "
	}

	if filter.ID > 0 {
		sqlQuery += " AND id = " + fmt.Sprint(filter.ID) + " "
	}

	sqlQuery = fmt.Sprintf("SELECT * FROM url WHERE 1=1 %s LIMIT 1;", sqlQuery)

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
		err = res.Scan(&id, &alias, &url)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}
		return &urlshort.URLShort{
			ID:    id,
			URL:   url,
			Alias: alias,
		}, nil
	}
	return nil, errors.New("url not found")
}
