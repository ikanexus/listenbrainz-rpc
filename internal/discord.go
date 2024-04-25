package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/hugolgst/rich-go/client"
	"github.com/spf13/viper"
)

type ScrobbleActivity struct {
	Album     string
	Artist    string
	Track     string
	Cover     string
	ReleaseId string
}

type DiscordActivity interface {
	Login()
	Logout()
	AddActivity(musicActivity *ScrobbleActivity) error
}

type discordActivity struct {
	active           bool
	appId            string
	listenbrainzUser string
}

func NewDiscordActivity() DiscordActivity {
	appId := viper.GetString("app-id")
	user := viper.GetString("user")
	return &discordActivity{
		appId:            appId,
		listenbrainzUser: user,
	}
}

func (d *discordActivity) Login() {
	log.Debugf("Logging in with AppID: %s", d.appId)
	if d.active {
		log.Debug("Already logged in, skipping")
		return
	}
	err := client.Login(d.appId)
	if err != nil {
		log.Fatalf("Unable to login to Discord IPC => %v", err)
	}
	d.active = true
}

func (d *discordActivity) Logout() {
	log.Debugf("Logging out from IPC")
	if d.active {
		client.Logout()
	} else {
		log.Debug("Already logged out, skipping")
	}
}

func (d discordActivity) getButtons(musicActivity *ScrobbleActivity) []*client.Button {
	// TODO: what do I do if there is no release ID
	return []*client.Button{
		{Label: "ListenBrainz Profile", Url: fmt.Sprintf("https://listenbrainz.org/user/%s", d.listenbrainzUser)},
		{Label: "Open on MusicBrainz", Url: fmt.Sprintf("https://musicbrainz.org/release/%s", musicActivity.ReleaseId)},
		// {Label: "Open on Spotify", Url: fmt.Sprintf("https://open.spotify.com/%s", "TODO")},
	}
}

func (d discordActivity) AddActivity(musicActivity *ScrobbleActivity) error {
	startTime := time.Now()
	timestamps := &client.Timestamps{Start: &startTime}
	activity := client.Activity{
		State:      musicActivity.Artist,
		Details:    musicActivity.Track,
		LargeImage: musicActivity.Cover,
		LargeText:  musicActivity.Album,
		// TODO: make small image optional / toggleable
		SmallImage: LISTENBRAINZ_LOGO,
		SmallText:  "Scrobbling via ListenBrainz",
		// TODO: can I set it to automatically expire when the song length is reached rather than waiting for LB to reset
		Timestamps: timestamps,
		Buttons:    d.getButtons(musicActivity),
	}
	return client.SetActivity(activity)
}
