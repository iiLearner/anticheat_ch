package auth

import (
	"anticheat_ch/data"
	"anticheat_ch/utils"
	"anticheat_ch/vars"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func FinalAuth(authuser string, db *sql.DB, key string, discord *discordgo.Session){

	if authuser == "open" {
		_, i, i2, i3, i4, i5, _ := data.GetTourneyInfo(db, key)
		vars.TournamentName = i
		tserver, _ := discord.Guild(i3)

		vars.TournamentServer = tserver.Name
		tuser, _ := discord.User(i2)
		vars.TournamentOrganizer = tuser.Username
		vars.ChannelID = i4
		vars.AlertsChannel = i5

		_, UID, userStatus := data.GetPlayerInfo(db, key)
		userObject, _ := discord.User(UID)
		vars.UserName = userObject.Username
		vars.UserID = UID;

		if userStatus == "-1"{
			fmt.Println("[ERROR] Your are banned from this tournament/software!")
			utils.CloseTerminal()
		}
	}else if authuser == "closed"{
		fmt.Println("[ERROR] The tournament is closed!")
		utils.CloseTerminal()
	}else{
		fmt.Println("[ERROR] Authenticationn failed! Are you sure that's the right key?")
		utils.CloseTerminal()
	}
}