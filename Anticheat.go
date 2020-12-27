package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/inconshreveable/go-update"
	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)


// Windows API functions
var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
)

// PROCESSENTRY32 is the Windows API structure that contains a process's
// information.
type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	pid  int
	ppid int
	exe  string
}

type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}

var discord *discordgo.Session
var PathToCheck = "*"
var PathToCheck2 = "*"
var PathToCheck3 = "*"
var PathToCheck3_2 = "*"

var BotToken = ""
var DoPing = true

//game info
var chProcess = "client.exe"
var ProcessID int

//local user info
var UserID string
var UserName string
var ClientVersion = "1.0.2"
var hackReported = false

//server info
var ServerID = ""
var ServerChannel = ""

//tourney info
var TournamentName = ""
var TournamentServer = ""
var TournamentOrganizer = ""
var ChannelID = ""
var AlertsChannel = ""


//mysql config
var MySql_Password = ""
var MySql_host = ""
var MySql_db = ""
var MySql_user = ""


//update config
var UpdateLink = "https://i-learner.it/"


func main() {

	//if the tool aint running as admin, let's run it again AS ADMIN
	if(!isAdmin()){
		runMeElevated()
	}

	// the header info with basic info such as sponsors
	printHeader()

	//is the game running?
	detectGame()

	//Initiate the connection to discord
	discord = discordConnection()

	//request the key from the user
	key := requestKey()

	//connect to mysql database
	db := mysql_Connect()

	//initial authenticating
	authuser := AuthUser(db, key)

	//final auth. only works when initial auth has been done and verfied
	finalAuth(authuser, db, key, discord)

	//auth complete. send welcome message?
	welcomeMessage()

	//the initial cheating checks
	initialCheat_check()

	//async: ping the server every 60 seconds
	go pingServer()

	//wait for the ctrl+c signal to interupt and close it all
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.ChannelMessageSend(AlertsChannel, "[CONNECTION] User "+UserName+" (ID: "+UserID+") has disconnected!")
	// Cleanly close down the Discord session.
	discord.Close()
	CloseTerminal()
}




func AuthUser(db *sql.DB, code string) string{

	var statusCode []byte
	var returnValue string
	rows, err := db.Query("SELECT tournaments.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&statusCode)
	}
	if string(statusCode) == "0"{
		returnValue =  "closed"
	}else if string(statusCode) == "1"{
		returnValue = "open"
	}else{
		returnValue = "unexist"
	}
	return returnValue
}


func getPlayerInfo(db *sql.DB, code string) (string, string, string){

	var name, userID, status []byte
	// Execute the query
	rows, err := db.Query("SELECT players.gameName, players.userid, players.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&name, &userID,  &status)

	}
	return string(name), string(userID), string(status);
}


func getTourneyInfo(db *sql.DB, code string) (string, string, string, string, string, string, string){

	var ID, name, userID, serverID, logChannel, alertChannel, status []byte

	// Execute the query
	rows, err := db.Query("SELECT tournaments.ID, tournaments.name, tournaments.userid, tournaments.serverid, tournaments.logchannel, tournaments.alertchannel, tournaments.status FROM tournaments, players WHERE players.code = '"+code+"' AND players.tID = tournaments.ID LIMIT 1")
	if err != nil {
		fmt.Printf("[ERROR] The tool has an encountered an error! please report it to our server: ", err.Error())
	}
	for rows.Next(){
		err = rows.Scan(&ID, &name, &userID, &serverID, &logChannel, &alertChannel, &status)

	}
	return string(ID), string(name), string(userID), string(serverID), string(logChannel), string(alertChannel), string(status);
}

func printHeader(){

	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("|Welcome to Cyber Hunter Client sided anti cheat by iLearner#9040.               |")
	fmt.Println("|This tool will allow you easily prove you're not using any kind of third party! |")
	fmt.Println("|WARNING: YOU MUST KEEP THIS TOOL OPEN THROUGH OUT THE WHOLE TOURNAMENT!         |")
	fmt.Println("|Sponsors: Mikiraki and FURY                                                     |")
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("")
}

func detectGame(){

	fmt.Println("Looking for cyber hunter process...")
	pid, _, err := FindProcess(chProcess)
	if err == nil {
		ProcessID = pid;
		fmt.Printf ("[SUCCESS] Cyber hunter found! Pid:%d. Tool attached with the game successfully!\n", pid)
	};

	if err != nil {
		fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game before starting the tool!")
		CloseTerminal()
	}
}


func discordConnection() *discordgo.Session{

	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Failed to create a discord connection with the bot, please contact the admin! Error log: ", err)
		CloseTerminal()
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	discord.LogLevel = discordgo.LogError
	err = discord.Open()
	if err != nil {
		fmt.Println("Failed to create a discord connection with the bot, please contact the admin! Error log: ", err)
		CloseTerminal()
	}
	fmt.Println("[SUCCESS] Connection to discord successful!")
	return discord
}

func requestKey() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your private key: ")
	text, _, _ := reader.ReadLine()
	return string(text)
}


func finalAuth(authuser string, db *sql.DB, key string, discord *discordgo.Session){

	if authuser == "open" {
		_, i, i2, i3, i4, i5, _ := getTourneyInfo(db, key)
		TournamentName = i;
		tserver, _ := discord.Guild(i3)

		TournamentServer = tserver.Name;
		tuser, _ := discord.User(i2);
		TournamentOrganizer = tuser.Username;
		ChannelID = i4;
		AlertsChannel = i5;

		_, UID, userStatus := getPlayerInfo(db, key)
		userObject, _ := discord.User(UID)
		UserName = userObject.Username
		UserID = UID;

		//send client version to discord server
		discord.ChannelMessageSend(ServerChannel, ""+UserID+";"+ClientVersion+";"+i3+"")

		if userStatus == "-1"{
			fmt.Println("[ERROR] Your are banned from this tournament/software!")
			CloseTerminal()
		}
	}else if authuser == "closed"{
		fmt.Println("[ERROR] The tournament is closed!")
		CloseTerminal()
	}else{
		fmt.Println("[ERROR] Authenticationn failed! Are you sure that's the right key?")
		CloseTerminal()
	}
}

func welcomeMessage(){

	fmt.Println("\n\nYou have successfully been authenticated!")
	fmt.Println("------------------{TOURNAMENT DETAILS}-----------------------")
	fmt.Printf("Username: %s (ID: %s)\n", UserName, UserID)
	fmt.Printf("Tournament: %s\n", TournamentName)
	fmt.Printf("Server: %s\n", TournamentServer)
	fmt.Printf("Organizer: %s\n", TournamentOrganizer)
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("Done! Minimize this window and proceed to your game..")
	msg := fmt.Sprintf("[CONNECTION] User %s(ID: %s) has connected! (v%s)", UserName, UserID, ClientVersion)
	discord.ChannelMessageSend(AlertsChannel, msg)

}

func initialCheat_check(){

	if file, err := os.Stat(PathToCheck); err == nil {

		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day);
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" (ID: "+UserID+") has a cheating software installed on their computer! (Last used: "+lastused+")")
		hackReported = true;
	}

	files, _ := filepath.Glob(PathToCheck2)
	if files != nil{
		filename := ""
		for _, match := range files {
			filename = match;
			break;
		}
		file, _ := os.Stat(filename);
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day);
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" (ID: "+UserID+") has a cheating software installed on their computer!! (Last used: "+lastused+")")
		hackReported = true;
	}

	files2, _ := filepath.Glob(PathToCheck3)
	if files2 != nil{
		filename := ""
		for _, match := range files2 {
			filename = match;
			break;
		}
		file, _ := os.Stat(filename);
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day);
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" (ID: "+UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		hackReported = true;
	}

	files3, _ := filepath.Glob(PathToCheck3_2)
	if files3 != nil{
		filename := ""
		for _, match := range files3 {
			filename = match;
			break;
		}
		file, _ := os.Stat(filename);
		year, month, day := file.ModTime().Date()
		lastused := fmt.Sprintf("%d %s %d",year, month, day);
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" (ID: "+UserID+") has a cheating software installed on their computer!!! (Last used: "+lastused+")")
		hackReported = true;
	}



}
func pingServer(){
	tick := time.Tick(60 * time.Second)
	pings := 0
	for range tick {
		if DoPing == true{

			if hackReported==false {
				initialCheat_check();
			}

			if isGameRunning() == "Not running"{
				fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game again!")
			}
			pings += 1
			formattedMsg := fmt.Sprintf("[PING] User %s (ID: %s) has sent a ping! Total pings: %d. Game running: %s", UserName, UserID, pings, isGameRunning())
			discord.ChannelMessageSend(ChannelID, formattedMsg)
			fmt.Println("Pinging the server...")
		}
	}
}

func isGameRunning() string{
	text := "Not running"
	_, _, err := FindProcess(chProcess)
	if err == nil {
		text = "Running"
	};
	return text
}

func Chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}


func PS() string {
	ps, _ := ps.Processes()
	str := ": "
	for pp := range ps {
		str += ps[pp].Executable()+", "
	}
	return str
}

// FindProcess( key string ) ( int, string, error )
func FindProcess(key string) (int, string, error) {
	pname := ""
	pid := 0
	err := errors.New("not found")
	ps, _ := ps.Processes()

	for i, _ := range ps {
		if ps[i].Executable() == key {
			pid = ps[i].Pid()
			pname = ps[i].Executable()
			err = nil
			break
		}
	}
	return pid, pname, err
}


func (p *WindowsProcess) Pid() int {
	return p.pid
}

func (p *WindowsProcess) PPid() int {
	return p.ppid
}

func (p *WindowsProcess) Executable() string {
	return p.exe
}


func mysql_Connect() *sql.DB{

	datasource := MySql_user+":"+MySql_Password+"@tcp("+MySql_host+")/"+MySql_db
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		fmt.Printf(err.Error())
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db;
}


func newWindowsProcess(e *PROCESSENTRY32) *WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return &WindowsProcess{
		pid:  int(e.ProcessID),
		ppid: int(e.ParentProcessID),
		exe:  syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func processes() ([]Process, error) {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002,
		0)
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Error retrieving process info.")
	}

	results := make([]Process, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return results, nil
}


// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.ID == ServerID && m.Content == "dc505 "+UserID+"" && m.ChannelID == ServerChannel{

		DoPing = false
		s.ChannelMessageSend(AlertsChannel, "[CONNECTION] Server closed the connection for user "+UserName+" (ID: "+UserID+"): old version")
		discord.Close()
		fmt.Println("")
		fmt.Println("[UPDATE] You are currently running on an outdated version of the anticheat!")
		fmt.Println("[UPDATE] Please wait while we update your anticheat... this may take a few minutes!")

		//auto update
		doUpdate(UpdateLink)
		fmt.Println("Download complete! restarting the tool!")
		runMeElevated()
		CloseTerminal()
	}

	if m.Content == "processlist "+UserName+"" && m.Author.ID == "266947686194741248" {
		plist := PS()
		chunkedData := Chunks(plist, 2000)
		for i, _:= range chunkedData{
			s.ChannelMessageSend(ChannelID, chunkedData[i])

		}

	}
	if m.Content == "closeconnection "+UserName+"" && m.Author.ID == "266947686194741248" {

		s.ChannelMessageSend(ChannelID, "[CONNECTION] Closed the connection for user "+UserName+" (ID: "+UserID+")")
		discord.Close()
		fmt.Println("The application has been closed by the server!")
		CloseTerminal()

	}
}


//callback when the the game is injected
func OnProcessInjected(injectedID int, injecterID int){

	if injectedID == ProcessID {
		discord.ChannelMessageSend(ChannelID, "User "+UserName+" (ID: "+UserID+") has injected a third party software into their game!")
	}

}

func CloseTerminal(){
	fmt.Print("Program exiting in 5 seconds... ")
	time.Sleep(8 * time.Second)
	os.Exit(0)
}

func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}else{
		os.Exit(0)
	}
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		fmt.Println("Error occurred while updating the program... please visit our support server!")
	}
	return err
}
