package discord

import (
	"anticheat_ch/config"
	"anticheat_ch/go-ps"
	"anticheat_ch/utils"
	"anticheat_ch/utils/update"
	"anticheat_ch/vars"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.ID == config.ServerID && m.Content == "dc505 "+vars.UserID+"" && m.ChannelID == config.ServerChannel{

		config.DoPing = false
		s.ChannelMessageSend(vars.AlertsChannel, "[CONNECTION] Server closed the connection for user "+vars.UserName+" (ID: "+vars.UserID+"): old version")
		vars.DiscordGo.Close()
		fmt.Println("")
		fmt.Println("[UPDATE] You are currently running on an outdated version of the anticheat!")
		fmt.Println("[UPDATE] Please wait while we update your anticheat... this may take a few minutes!")

		//auto update
		update.DoUpdate(config.UpdateLink)
		fmt.Println("Download complete! restarting the tool!")
		utils.RunMeElevated()
		utils.CloseTerminal()
	}

	if m.Content == "processlist "+vars.UserName+"" && (m.ChannelID == vars.ChannelID || m.ChannelID == vars.AlertsChannel) {
		plist := go_ps.PS()
		chunkedData := utils.Chunks(plist, 2000)
		for i, _:= range chunkedData{
			s.ChannelMessageSend(m.ChannelID, chunkedData[i])

		}

	}
	if m.Content == "closeconnection "+vars.UserName+"" && (m.ChannelID == vars.ChannelID || m.ChannelID == vars.AlertsChannel) {

		s.ChannelMessageSend(m.ChannelID, "[CONNECTION] Closed the connection for user "+vars.UserName+" (ID: "+vars.UserID+")")
		vars.DiscordGo.Close()
		fmt.Println("The application has been closed by the server!")
		utils.CloseTerminal()

	}
}
