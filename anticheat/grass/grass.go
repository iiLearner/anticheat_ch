package grass

import (
	"anticheat_ch/vars"
	"anticheat_ch/windows"
	"fmt"
	"os"
	"strings"
	"time"
)

func GrassCheck(){
	path := windows.Path(vars.ProcessID)
	path = path+"/Documents/res/levelsets/g80.npk"
	path = strings.Replace(path, "\\bin\\client.exe", "", -1)
	if file, err := os.Stat(path); err == nil {

		tdr := file.ModTime().Unix()
		tdrn := time.Now().Unix()
		minsago := (tdr-tdrn)/60

		if (tdrn-tdr)<3600 {
			msg := fmt.Sprintf("[N0 GRASS WARNING] User %s(ID: %s) could POTENTIALLY BE using NO-GRASS! (Files modified %d mins ago)", vars.UserName, vars.UserID, minsago)
			vars.DiscordGo.ChannelMessageSend(vars.AlertsChannel, msg)
		}
	}
}

