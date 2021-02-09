package anticheat

import (
	"anticheat_ch/config"
	"anticheat_ch/utils/game"
	"anticheat_ch/vars"
	"bufio"
	"fmt"
	"os"
	"time"
)

func RequestKey() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your private key: ")
	text, _, _ := reader.ReadLine()
	return string(text)
}

func PingServer(){
	tick := time.Tick(60 * time.Second)
	pings := 0
	for range tick {
		if config.DoPing == true{

			if vars.HackReported==false {
				InitialCheat_check()
			}

			if game.IsGameRunning() == "Not running"{
				fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game again!")
			}
			pings += 1
			formattedMsg := fmt.Sprintf("[PING] User %s (ID: %s) has sent a ping! Total pings: %d. Game running: %s", vars.UserName, vars.UserID, pings, game.IsGameRunning())
			vars.DiscordGo.ChannelMessageSend(vars.ChannelID, formattedMsg)
			fmt.Println("Pinging the server...")
		}
	}
}

