package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


func initializeMysql(database Database) (*sql.DB, error) {
	d, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&collation=utf8_unicode_ci&parseTime=True",
			database.Username,
			database.Password,
			database.Host,
			database.Port,
			database.DBName,
		),
	)
	if err != nil {
		log.Panicln(err)
	}

	d.SetMaxOpenConns(database.MaxOpenConnections)
	d.SetMaxIdleConns(database.MaxIdleConnections)

	if err := d.Ping(); err != nil {
		return nil, err
	}

	return d, nil
}

