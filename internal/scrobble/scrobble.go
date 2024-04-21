package scrobble

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ikanexus/listenbrainz-rpc/internal/coverartarchive"
	"github.com/ikanexus/listenbrainz-rpc/internal/discord"
	"github.com/ikanexus/listenbrainz-rpc/internal/listenbrainz"
)

const SLEEP_DELAY = 10 * time.Second

func Scrobble(username string) error {
	log.Printf("Running with user: %s\n", username)

	var currentSongId string

	for running := true; running; {
		currentSong, err := checkCurrentPlaying(username)

		if err != nil {
			return err
		}

		trackHash := getTrackHash(currentSong.TrackMetadata)
		if trackHash == currentSongId {
			log.Printf("Still playing %s", currentSong.TrackMetadata.TrackName)
			time.Sleep(SLEEP_DELAY)
			continue
		}
		currentSongId = trackHash

		trackMetadata := currentSong.TrackMetadata
		trackName := trackMetadata.TrackName
		artistName := trackMetadata.ArtistName
		albumName := trackMetadata.ReleaseName

		releaseId := getReleaseId(trackMetadata)
		albumArt := coverartarchive.GetAlbumArt(releaseId)

		activity := &discord.ScrobbleActivity{
			Album:     albumName,
			Artist:    artistName,
			Track:     trackName,
			Cover:     albumArt,
			ReleaseId: releaseId,
		}

		err = discord.AddMusicActivity(activity)

		if err != nil {
			return err
		}
	}
	return nil
}

func checkCurrentPlaying(username string) (*listenbrainz.Listen, error) {
	nowPlaying, err := listenbrainz.GetUserNowPlaying(username)
	if err != nil {
		return nil, err
	}
	currentSong := listenbrainz.FindCurrentListen(*nowPlaying)

	if currentSong == nil {
		fmt.Println("No song playing")
		return nil, nil
	}
	return currentSong, nil
}

func getDuration(trackInfo listenbrainz.TrackAdditionalInfo) int {
	if trackInfo.Duration != 0 {
		return trackInfo.Duration
	}
	if trackInfo.DurationMs != 0 {
		return trackInfo.DurationMs / 1000
	}
	return 0
}

func getReleaseId(metadata listenbrainz.TrackMetadata) string {
	locations := []listenbrainz.MbidMapping{
		metadata.AdditionalInfo.MbidMapping,
		metadata.MbidMapping,
	}
	var result string

	for _, location := range locations {
		releaseId := location.ReleaseMBID
		if releaseId != "" {
			result = releaseId
			break
		}
	}
	return result
}

func getTrackHash(metadata listenbrainz.TrackMetadata) string {
	b, _ := json.Marshal(&metadata)
	return string(b)
}

func getTrackId(metadata listenbrainz.TrackMetadata) string {
	locations := []listenbrainz.MbidMapping{
		metadata.AdditionalInfo.MbidMapping,
		metadata.MbidMapping,
	}
	var result string

	for _, location := range locations {
		trackId := location.TrackMBID
		if trackId != "" {
			result = trackId
			break
		}
	}
	return result
}
