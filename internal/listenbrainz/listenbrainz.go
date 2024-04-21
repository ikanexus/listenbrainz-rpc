package listenbrainz

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const API_VERSION = 1
const API_BASE = "https://api.listenbrainz.org"

func GetUserNowPlaying(user string) (*ListenPayload, error) {
	res, err := http.Get(fmt.Sprintf("%s/%d/user/%s/playing-now", API_BASE, API_VERSION, user))
	if err != nil {
		return nil, err
	}

	body := res.Body
	defer body.Close()

	var listenResponse ListenResponse

	err = json.NewDecoder(body).Decode(&listenResponse)
	if err != nil {
		return nil, err
	}

	return &listenResponse.Payload, nil

}

func FindCurrentListen(payload ListenPayload) *Listen {
	var currentListen *Listen
	if payload.PlayingNow == false {
		return nil
	}
	if payload.Count == 0 {
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
