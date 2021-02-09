package data

import (
	"database/sql"
	"fmt"
)

func GetTourneyInfo(db *sql.DB, code string) (string, string, string, string, string, string, string){

	var ID, name, userID, serverID, logChannel, alertChannel, status []byte

	// Execute the query
	rows, err := db.Query("SELECT tournaments.ID, tournaments.name, tournaments.userid, tournaments.serverid, tournaments.logchannel, tournaments.alertchannel, tournaments.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&ID, &name, &userID, &serverID, &logChannel, &alertChannel, &status)

	}
	return string(ID), string(name), string(userID), string(serverID), string(logChannel), string(alertChannel), string(status);
}
