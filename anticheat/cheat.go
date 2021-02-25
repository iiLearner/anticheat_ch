package anticheat

import (
	"anticheat_ch/config"
	"anticheat_ch/vars"
	"fmt"
	"os"
)

func InitialcheatCheck(){

	if file, err := os.Stat(vars.VEnv1); err == nil {

		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer! (Venom cheats)")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer! (Last used: "+lastused+")")
		vars.HackReported = true
	}

	homedir, _ := os.UserHomeDir()
	dwnPath := homedir + "/Downloads/*"
	dstp :=  homedir + "/Desktop/"

	CheckCheats(dwnPath+"pvip*", "PVIP Hacks")
	CheckCheats(dstp+"pvip*", "PVIP Hacks")
	CheckCheats(vars.VEnv32, "PVIP Hacks")
	CheckCheats(vars.VEnv3, "PVIP Hacks")

	CheckCheats(dwnPath+"HIBARAGAMER*", "Hibra Hacks")
	CheckCheats(dstp+"HIBARAGAMER*", "Hibra Hacks")
	CheckCheats(vars.VEnv2, "Hibra Hacks")


}
