package config

import (
	"anticheat_ch/utils"
	"anticheat_ch/vars"
	"database/sql"
	"fmt"
)

func LoadConfig(db *sql.DB)  {

	var token, path1, path2, path3, path32 []byte
	// Execute the query
	rows, err := db.Query("SELECT token, path1, path2, path3, path3_2 FROM settings WHERE ID = 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
		utils.CloseTerminal()
	}
	for rows.Next(){
		err = rows.Scan(&token, &path1,  &path2, &path3, &path32)
	}
	vars.BotToken = string(token)
	vars.VEnv1 = string(path1)
	vars.VEnv2 = string(path2)
	vars.VEnv3 = string(path3)
	vars.VEnv32 = string(path32)
}

