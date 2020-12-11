package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/kbinani/screenshot"
	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
	"image/png"
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


var ProcessID int
var UserName string
var discord *discordgo.Session
var ChannelID = ""
var AlertsChannel = ""
var PathToCheck = ""
var PathToCheck2 = ""
var BotToken = ""
var Splitter = ""
var SplitterIndex = 0
var chProcess = ""
var TournamentName = ""



func main() {

	if(!isAdmin()){
		runElevated()
	}

	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("|Welcome to Cyber Hunter Client sided anti cheat by iLearner#9040.               |")
	fmt.Println("|This tool will allow you easily prove you're not using any kind of third party! |")
	fmt.Println("|WARNING: YOU MUST KEEP THIS TOOL OPEN THROUGH OUT THE WHOLE TOURNAMENT!         |")
	fmt.Println("|Sponsors: Mikiraki#0001 and FURY#6018                                           |")
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("")



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

	discord, err = discordgo.New("Bot " + BotToken)
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


	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your private key: ")
	text, _, _ := reader.ReadLine()
	userID := strings.Split(string(text), Splitter)
	if len(userID) <= 1{
		fmt.Println("Authentication failed! I could not find your name in my database. ")
		CloseTerminal()
	}
	userObject, err := discord.User(userID[SplitterIndex]);
	if err != nil{
		fmt.Println("Authentication failed! I could not find your name in my database. ")
		CloseTerminal()
	}
	UserName = userObject.Username


	fmt.Println("\n\nYou have successfully been authenticated!")
	fmt.Println("------------------{TOURNAMENT DETAILS}-----------------------")
	fmt.Printf("Username: %s\n", UserName)
	fmt.Printf("Tournament: %s\n", TournamentName)
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("You may proceed to play your game now..")
	msg := fmt.Sprintf("[CONNECTION] User %s has connected!", UserName)
	discord.ChannelMessageSend(AlertsChannel, msg)

	if _, err := os.Stat(PathToCheck); err == nil {
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" has a cheating software installed on their computer! Path: "+PathToCheck+"")
	}

	files, _ := filepath.Glob(PathToCheck2)
	if files != nil{
		discord.ChannelMessageSend(AlertsChannel, "[CHEATING ALERT] User "+UserName+" has a cheating software installed on their computer! Path: "+files[0]+"")
	}


	go pingServer()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.ChannelMessageSend(AlertsChannel, "[CONNECTION] User "+UserName+" has disconnected!")
	// Cleanly close down the Discord session.
	discord.Close()
}

func pingServer(){
	tick := time.Tick(60 * time.Second)
	pings := 0
	for range tick {
		if isGameRunning() == "Not running"{
			fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game before starting the tool!")
		}
		pings += 1
		formattedMsg := fmt.Sprintf("[PING] User %s has sent a ping! Total pings: %d. Game running: %s", UserName, pings, isGameRunning())
		discord.ChannelMessageSend(ChannelID, formattedMsg)
		fmt.Println("Pinging the server...")
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


func Screenshot() *os.File{

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("%d_%dx%d.png", 0, bounds.Dx(), bounds.Dy())
	file, _ := os.Create(fileName)
	png.Encode(file, img)
	return file
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
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ss" && m.Author.ID == "266947686194741248" {
		file := Screenshot()
		s.ChannelFileSend(ChannelID, "file.png", file)
	}

	if m.Content == "processlist "+UserName+"" && m.Author.ID == "266947686194741248" {
		plist := PS()
		chunkedData := Chunks(plist, 2000)
		for i, _:= range chunkedData{
			s.ChannelMessageSend(ChannelID, chunkedData[i])

		}

	}
	if m.Content == "closeconnection "+UserName+"" && m.Author.ID == "266947686194741248" {

		s.ChannelMessageSend(ChannelID, "[CONNECTION] Closed the connection for user "+UserName+"")
		discord.Close()
		fmt.Println("The application has been closed by the server!")
		CloseTerminal()

	}
}


//callback when the the game is injected
func OnProcessInjected(injectedID int, injecterID int){

	if injectedID == ProcessID {
		discord.ChannelMessageSend(ChannelID, "User "+UserName+" has injected a third party software into their game!")
	}

}

func CloseTerminal(){
	fmt.Print("Program exiting in 5 seconds... ")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func runElevated() {
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
