package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func MySQL_Connect() *sql.DB{


	datasource := os.Getenv("MySql_user")+":"+os.Getenv("MySql_Password")+"@tcp("+os.Getenv("MySql_host")+")/"+os.Getenv("MySql_db")
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		fmt.Printf(err.Error())
	}
	if err == nil{
		fmt.Println("[SUCCESS] Connection to database was successful!")
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
