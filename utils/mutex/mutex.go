package mutex

import (
	"anticheat_ch/utils"
	"anticheat_ch/windows"
	"fmt"
	"syscall"
	"unsafe"
)

func MutexCheck(){

	_, err := CreateMutex("iAC")
	if err!=nil{
		fmt.Printf("Another instance of the program is already running...")
		utils.CloseTerminal()
	}
}

func CreateMutex(name string) (uintptr, error) {
	ret, _, err := windows.ProcCreateMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
	)
	switch int(err.(syscall.Errno)) {
	case 0:
		return ret, nil
	default:
		return ret, err
	}
}
