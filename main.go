package main

import (
	"anticheat_ch/anticheat"
	"anticheat_ch/auth"
	"anticheat_ch/config"
	"anticheat_ch/discord"
	"anticheat_ch/message"
	"anticheat_ch/utils"
	"anticheat_ch/utils/game"
	"anticheat_ch/utils/mutex"
	"anticheat_ch/utils/mysql"
	"anticheat_ch/utils/update"
	"anticheat_ch/vars"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
	"syscall"
)


func main() {

	//if the tool aint running as admin, let's run it again AS ADMIN
	if !utils.IsAdmin() {
		utils.RunMeElevated()
	}

	mutex.MutexCheck()

	// the header info with basic info such as sponsors
	message.PrintHeader()

	//is the game running?
	game.DetectGame()

	//connect to mysql database
	db := mysql.MySQL_Connect()

	//load config from db
	config.LoadConfig(db)

	//Initiate the connection to discord
	vars.DiscordGo = discord.DiscordConnection()

	//check update
	update.CheckUpdate()

	//request the key from the user
	key := anticheat.RequestKey()

	//initial authenticating
	authuser := auth.AuthUser(db, key)

	//final auth. only works when initial auth has been done and verfied
	auth.FinalAuth(authuser, db, key, vars.DiscordGo)

	//user must be in support server
	check := auth.DiscordCheck()

	//auth complete. send welcome message?
	if check {message.WelcomeMessage()}

	//check if user has no grass
	//grass.GrassCheck()

	//the initial cheating checks
	anticheat.InitialcheatCheck()

	//async: ping the server every 60 seconds
	go anticheat.PingServer()

	//wait for the ctrl+c signal to interupt and close it all
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, "[CONNECTION] User "+vars.UserName+" (ID: "+vars.UserID+") has disconnected!")
	// Cleanly close down the Discord session.
	vars.DiscordGo.Close()
	utils.CloseTerminal()
}
