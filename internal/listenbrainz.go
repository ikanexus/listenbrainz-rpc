package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

const API_VERSION = 1
const API_BASE = "https://api.listenbrainz.org"

type ListenResponse struct {
	Payload ListenPayload `json:"payload"`
}

type ListenPayload struct {
	Count      int      `json:"count"`
	PlayingNow bool     `json:"playing_now"`
	UserId     string   `json:"user_id"`
	Listens    []Listen `json:"listens"`
}

type Listen struct {
	PlayingNow    bool          `json:"playing_now"`
	ListenedAt    int           `json:"listened_at"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
	AdditionalInfo TrackAdditionalInfo `json:"additional_info"`
	ArtistName     string              `json:"artist_name"`
	TrackName      string              `json:"track_name"`
	ReleaseName    string              `json:"release_name"`
	MbidMapping    MbidMapping         `json:"mbid_mapping"`
}

type TrackAdditionalInfo struct {
	MbidMapping
	TrackNumber int      `json:"track_number"`
	Tags        []string `json:"tags"`
	Duration    int      `json:"duration"`
	DurationMs  int      `json:"duration_ms"`
}

type MbidMapping struct {
	ReleaseMBID      string   `json:"release_mbid"`
	ArtistMBIDs      []string `json:"artist_mbids"`
	ReleaseGroupMBID string   `json:"release_group_mbid"`
	RecordingMBID    string   `json:"recording_mbid"`
	TrackMBID        string   `json:"track_mbid"`
}

type listenBrainz struct {
	username string
}

type ListenBrainz interface {
	GetNowPlaying() *Listen
}

func (l listenBrainz) GetNowPlaying() *Listen {
	res, err := http.Get(fmt.Sprintf("%s/%d/user/%s/playing-now", API_BASE, API_VERSION, l.username))
	if err != nil {
		log.Errorf("Unable to get current playing => %v", err)
		return nil
	}

	body := res.Body
	defer body.Close()

	var listenResponse ListenResponse

	err = json.NewDecoder(body).Decode(&listenResponse)
	if err != nil {
		log.Errorf("Unable to decode json body => %v", err)
		return nil
	}

	payload := listenResponse.Payload

	return l.findCurrentListen(&payload)

}

func (l listenBrainz) findCurrentListen(payload *ListenPayload) *Listen {
	var currentListen *Listen
	if payload.PlayingNow == false || payload.Count == 0 {
		log.Debugf("No song playing")
		return nil
	}
	for _, listen := range payload.Listens {
		if listen.PlayingNow == true {
			currentListen = &listen
			break
		}
	}
	return currentListen
}

func NewListenBrainz(username string) ListenBrainz {
	return &listenBrainz{
		username: username,
	}
}
