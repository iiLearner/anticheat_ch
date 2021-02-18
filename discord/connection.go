package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"anticheat_ch/vars"
	"anticheat_ch/utils"
)

func DiscordConnection() *discordgo.Session{

	discord, err := discordgo.New("Bot " + vars.BotToken)
	if err != nil {
		fmt.Println("Failed to create a discord connection with the bot, please contact the admin! Error log: ", err)
		utils.CloseTerminal()
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	discord.LogLevel = discordgo.LogError
	err = discord.Open()
	if err != nil {
		fmt.Println("Failed to create a discord connection with the bot, please contact the admin! Error log: ", err)
		utils.CloseTerminal()
	}
	fmt.Println("[SUCCESS] Connection to discord successful!")
	return discord
}
