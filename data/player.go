package data

import (
	"database/sql"
	"fmt"
)

func GetPlayerInfo(db *sql.DB, code string) (string, string, string){

	var name, userID, status []byte
	// Execute the query
	rows, err := db.Query("SELECT players.gameName, players.userid, players.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&name, &userID,  &status)

	}
	return string(name), string(userID), string(status);
}
