package discord

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/hugolgst/rich-go/ipc"
)

type ScrobbleActivity struct {
	Album     string
	Artist    string
	Track     string
	Cover     string
	ReleaseId string
}

func Login(appId string) error {
	return client.Login(appId)
}

func Logout() {
	client.Logout()
}

func AddMusicActivity(musicActivity *ScrobbleActivity) error {
	startTime := time.Now()
	timestamps := &client.Timestamps{Start: &startTime}
	activity := client.Activity{
		State:      musicActivity.Artist,
		Details:    musicActivity.Track,
		LargeImage: musicActivity.Cover,
		Timestamps: timestamps,
		// TODO: handle when there is no release ID - doesn't break, just has an invalid URL
		Buttons: []*client.Button{
			{Label: "Open on MusicBrainz", Url: fmt.Sprintf("https://musicbrainz.org/release/%s", musicActivity.ReleaseId)},
		},
	}
	return client.SetActivity(activity)
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func ClearMusicActivity() error {
	payload, err := json.Marshal(client.Frame{
		Cmd: "CLEAR_ACTIVITY",
		Args: client.Args{
			Pid: os.Getpid(),
		},
		Nonce: getNonce(),
	})
	if err != nil {
		return nil
	}
	out := ipc.Send(1, string(payload))
	log.Printf("ipc clear: %s\n", out)
	return nil
}
