package windows

import (
	"anticheat_ch/vars"
	"fmt"
	"syscall"
	"unsafe"
)

// Path returns path to process executable
func Path(pid int) (string) {
	processModules, _ := modules(pid)
	return processModules[0].path
}

func newWindowsModule(e *MODULEENTRY32) windowsModule {
	return windowsModule{
		name: ptrToString(e.SzModule[:]),
		path: ptrToString(e.SzExePath[:]),
	}
}

func ptrToString(c []uint16) string {
	i := 0
	for {
		if c[i] == 0 {
			return syscall.UTF16ToString(c[:i])
		}
		i++
	}
}

func modules(pid int) ([]windowsModule, error) {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000008, // TH32CS_SNAPMODULE
		uintptr(uint32(pid)))
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry MODULEENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procModule32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("Error retrieving module info")
	}

	results := make([]windowsModule, 0, 50)
	for {
		results = append(results, newWindowsModule(&entry))

		ret, _, _ := procModule32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return results, nil
}

//callback when the the game is injected
func OnProcessInjected(injectedID int, injecterID int){

	if injectedID == vars.ProcessID {
		vars.DiscordGo.ChannelMessageSend(vars.ChannelID, "User "+vars.UserName+" (ID: "+vars.UserID+") has injected a third party software into their game!")
	}

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

func (p *WindowsProcess) Pid() int {
	return p.pid
}

func (p *WindowsProcess) PPid() int {
	return p.ppid
}

func (p *WindowsProcess) Executable() string {
	return p.exe
}