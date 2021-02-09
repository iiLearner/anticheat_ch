package anticheat

import (
	"anticheat_ch/config"
	"anticheat_ch/vars"
	"fmt"
	"os"
	"path/filepath"
)

func InitialCheat_check(){

	if file, err := os.Stat(vars.VEnv1); err == nil {

		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer! (Last used: "+lastused+")")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer! (Last used: "+lastused+")")
		vars.HackReported = true
	}

	files, _ := filepath.Glob(vars.VEnv2)
	if files != nil{
		filename := ""
		for _, match := range files {
			filename = match;
			break
		}
		file, _ := os.Stat(filename);
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!! (Last used: "+lastused+")")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!! (Last used: "+lastused+")")
		vars.HackReported = true
	}

	files2, _ := filepath.Glob(vars.VEnv3)
	if files2 != nil{
		filename := ""
		for _, match := range files2 {
			filename = match
			break
		}
		file, _ := os.Stat(filename)
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		vars.HackReported = true
	}

	files3, _ := filepath.Glob(vars.VEnv32)
	if files3 != nil{
		filename := ""
		for _, match := range files3 {
			filename = match
			break
		}
		file, _ := os.Stat(filename)
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day)
		vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		vars.DiscordGo.ChannelMessageSend(config.ServerAlerts, "[CHEATING ALERT] User "+vars.UserName+" (ID: "+vars.UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		vars.HackReported = true
	}
}
