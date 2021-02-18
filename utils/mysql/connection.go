package mysql

import (
	"anticheat_ch/config"
	"database/sql"
	"fmt"
	"time"
)

func MySQL_Connect() *sql.DB {

	datasource := ""+config.MySql_user+":"+config.MySql_Password+"@tcp("+config.MySql_host+")/"+config.MySql_db+""
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		fmt.Printf(err.Error())
	}
	if err == nil {
		fmt.Println("[SUCCESS] Connection to database was successful!")
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
