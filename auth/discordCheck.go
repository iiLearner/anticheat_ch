package auth

import (
	"anticheat_ch/utils"
	"anticheat_ch/vars"
	"fmt"
)

func DiscordCheck()  bool{

	check := true
	_, err := vars.DiscordGo.GuildMember("786699054540914688", vars.UserID)

	if err != nil{
		fmt.Println("")
		fmt.Println("[ERROR] You must be in support server in order to login!")
		fmt.Println("[ERROR] Use the invite link sent to you by iBot to join.")
		check = false
		utils.CloseTerminal()
	}
	return check

}
