package internal

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

const LISTENBRAINZ_LOGO = "https://archive.org/download/listenbrainz-20190401-000403/ListenBrainz_Logo.png"

type CoverArtRetriever interface {
	GetFront(releaseMbid string) string
}

type coverArtRetriever struct {
	client *http.Client
}

func NewCoverArtRetriever() CoverArtRetriever {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &coverArtRetriever{
		client: client,
	}
}

func (c coverArtRetriever) GetFront(releaseMbid string) string {
	size := 250
	defaultAlbumArt := LISTENBRAINZ_LOGO
	archiveArtUrl := fmt.Sprintf("https://coverartarchive.org/release/%s/front-%d.jpg", releaseMbid, size)

	response, err := c.client.Get(archiveArtUrl)
	if err != nil {
		log.Errorf("Error retrieving album art %s => %v\n", archiveArtUrl, err)
		return defaultAlbumArt
	}

	url, err := response.Location()
	if err != nil {
		log.Warnf("No location response header - defaulting to placeholder (%s) => %v\n", defaultAlbumArt, err)
		return defaultAlbumArt
	}

	return url.String()
}
