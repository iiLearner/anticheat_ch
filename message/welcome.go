package message

import (
	"anticheat_ch/config"
	"anticheat_ch/vars"
	"fmt"
)

func WelcomeMessage(){

	fmt.Println("\n\nYou have successfully been authenticated!")
	fmt.Println("------------------{TOURNAMENT DETAILS}-----------------------")
	fmt.Printf("Username: %s (ID: %s)\n", vars.UserName, vars.UserID)
	fmt.Printf("Tournament: %s\n", vars.TournamentName)
	fmt.Printf("Server: %s\n", vars.TournamentServer)
	fmt.Printf("Organizer: %s\n", vars.TournamentOrganizer)
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("Done! Minimize this window and proceed to your game..")
	msg := fmt.Sprintf("[CONNECTION] User %s(ID: %s) has connected! (v%s)", vars.UserName, vars.UserID, config.ClientVersion)
	vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, msg)

}
