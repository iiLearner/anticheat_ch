package game

import (
	"anticheat_ch/go-ps"
	"anticheat_ch/utils"
	"anticheat_ch/vars"
	"fmt"
)

func DetectGame(){

	fmt.Println("Looking for cyber hunter process...")
	pid, _, err1 := go_ps.FindProcess(vars.ChProcess)
	if err1 == nil {
		vars.ProcessID = pid;
		fmt.Printf ("[SUCCESS] Cyber hunter found! Pid:%d. Tool attached with the game successfully!\n", pid)
	};

	if err1 != nil {
		fmt.Println("[ERROR] The tool was unable to find cyber hunter, please run the game before starting the tool!")
		utils.CloseTerminal()
	}
}

func IsGameRunning() string{
	text := "Not running"
	_, _, err := go_ps.FindProcess(vars.ChProcess)
	if err == nil {
		text = "Running"
	}
	return text
}