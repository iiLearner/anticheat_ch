package auth

import (
	"database/sql"
	"fmt"
)

func AuthUser(db *sql.DB, code string) string{

	var statusCode []byte
	var returnValue string
	rows, err := db.Query("SELECT tournaments.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&statusCode)
	}
	if string(statusCode) == "0"{
		returnValue =  "closed"
	}else if string(statusCode) == "1"{
		returnValue = "open"
	}else{
		returnValue = "unexist"
	}
	return returnValue
}
