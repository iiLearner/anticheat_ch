package update

import (
	"anticheat_ch/config"
	"anticheat_ch/vars"
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/twinj/uuid"
	"net/http"
)

func CheckUpdate()  {

	//send client version to discord server
	vars.UserID = uuid.NewV4().String()
	vars.DiscordGo.ChannelMessageSend(config.ServerChannel, ""+vars.UserID+";"+config.ClientVersion+"")
}

func DoUpdate(url string) error {
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

