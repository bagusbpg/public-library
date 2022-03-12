package util

import (
	_config "plain-go/public-library/app/config"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func GetDBInstance(config *_config.AppConfig) (*sql.DB, error) {
	if db == nil {
		driverName, dataSourceName := config.Database.Driver, config.Database.Connection
		initdb, err := sql.Open(driverName, dataSourceName)

		if err != nil {
			return nil, err
		}

		db = initdb
	}

	return db, nil
}
