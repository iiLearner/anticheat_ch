package anticheat

import (
	"anticheat_ch/anticheat/cheats"
	"anticheat_ch/config"
	"anticheat_ch/utils/game"
	"anticheat_ch/vars"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RequestKey() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your private key: ")
	text, _, _ := reader.ReadLine()
	return string(text)
}

func CheckCheats(env string, name string)  {

	files4, _ := filepath.Glob(env)
	if files4 != nil{
		filename := ""
		for _, match := range files4 {
			filename = match
			break
		}
		file, _ := os.Stat(filename)
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer! ("+name+"))")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		vars.HackReported = true
	}

	
}

func PingServer() {
	tick := time.Tick(60 * time.Second)
	pings := 0
	for range tick {
		if config.DoPing == true {

			if vars.HackReported == false {
				InitialcheatCheck()
			}

			if vars.RecoilReported == false {
				cheats.RecoilCheck()
			}

			if game.IsGameRunning() == "Not running" {
				fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game again!")
			}
			pings += 1
			formattedMsg := fmt.Sprintf("[PING] User %s (ID: %s) has sent a ping! Total pings: %d. Game running: %s", vars.UserName, vars.UserID, pings, game.IsGameRunning())
			vars.DiscordGo.ChannelMessageSend(vars.ChannelID, formattedMsg)
			fmt.Println("Pinging the server...")
		}
	}
}
