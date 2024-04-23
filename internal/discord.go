package internal

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
)

type ScrobbleActivity struct {
	Album     string
	Artist    string
	Track     string
	Cover     string
	ReleaseId string
}

type DiscordActivity interface {
	Login() error
	Logout()
	AddActivity(musicActivity *ScrobbleActivity) error
}

type discordActivity struct {
	appId string
}

func NewDiscordActivity(appId string) DiscordActivity {
	return &discordActivity{
		appId: appId,
	}
}

func (d discordActivity) Login() error {
	return client.Login(d.appId)
}

func (d discordActivity) Logout() {
	client.Logout()
}

func (d discordActivity) getButtons() {

}

func (d discordActivity) AddActivity(musicActivity *ScrobbleActivity) error {
	startTime := time.Now()
	timestamps := &client.Timestamps{Start: &startTime}
	activity := client.Activity{
		State:      musicActivity.Artist,
		Details:    musicActivity.Track,
		LargeImage: musicActivity.Cover,
		LargeText:  musicActivity.Album,
		SmallImage: LISTENBRAINZ_LOGO,
		Timestamps: timestamps,
		Buttons: []*client.Button{
			{Label: "Open on MusicBrainz", Url: fmt.Sprintf("https://musicbrainz.org/release/%s", musicActivity.ReleaseId)},
		},
	}
	return client.SetActivity(activity)
}
