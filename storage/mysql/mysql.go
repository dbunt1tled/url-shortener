package mysql

import (
	"sync"
	"time"

	"github.com/dbunt1tled/url-shortener/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	instance *Mysql    //nolint:gochecknoglobals // singleton
	dm       sync.Once //nolint:gochecknoglobals // singleton
)

type Mysql struct {
	db *sqlx.DB
}

func (db *Mysql) Close() error {
	return db.db.Close()
}

func (db *Mysql) GetDB() *sqlx.DB {
	return db.db
}

func GetInstance() *Mysql {
	dm.Do(func() {
		instance = connection()
	})
	return instance
}

func connection() *Mysql {
	cfg := config.LoadConfig()
	db, err := sqlx.Connect("mysql", cfg.DatabaseDSN)
	if err != nil {
		panic("Error init connection:" + err.Error())
	}
	db.SetMaxOpenConns(1)                   // TODO: need determinate
	db.SetMaxIdleConns(1)                   // TODO: need determinate
	db.SetConnMaxLifetime(5 * time.Minute)  //nolint:mnd // TODO: need determinate
	db.SetConnMaxIdleTime(10 * time.Minute) //nolint:mnd // TODO: need determinate
	err = db.Ping()
	if err != nil {
		panic("Error ping:" + err.Error())
	}

	return &Mysql{db: db}
}
