package internal

import (
	"github.com/charmbracelet/log"
	"time"
)

const SLEEP_DELAY = 10 * time.Second

type scrobbler struct {
	username        string
	discordActivity DiscordActivity
	coverRetriever  CoverArtRetriever
	listenBrainz    ListenBrainz
}

type Scrobbler interface {
	Scrobble() error
}

func (s scrobbler) Scrobble() error {
	log.Infof("Running with user: %s", s.username)
	defer s.discordActivity.Logout()
	s.discordActivity.Login()

	isIdle := false

	var currentSongId string

	for running := true; running; {
		currentSong := s.checkCurrentPlaying()

		if currentSong == nil {
			time.Sleep(SLEEP_DELAY * 2)
			isIdle = true
			s.discordActivity.Logout()
			continue
		}

		if isIdle {
			s.discordActivity.Login()
			isIdle = false
		}

		releaseInfo := NewReleaseInfoRetriever(currentSong.TrackMetadata)

		trackHash := releaseInfo.GetTrackHash()
		trackMetadata := currentSong.TrackMetadata
		trackName := trackMetadata.TrackName
		artistName := trackMetadata.ArtistName
		albumName := trackMetadata.ReleaseName

		if trackHash == currentSongId {
			log.Debugf("Still playing %s by %s (%s)", trackName, artistName, albumName)
			time.Sleep(SLEEP_DELAY)
			continue
		}
		log.Infof("Playing %s by %s (%s)", trackName, artistName, albumName)
		currentSongId = trackHash

		releaseId := releaseInfo.GetReleaseId()
		albumArt := s.coverRetriever.GetFront(releaseId)

		activity := &ScrobbleActivity{
			Album:     albumName,
			Artist:    artistName,
			Track:     trackName,
			Cover:     albumArt,
			ReleaseId: releaseId,
		}

		err := s.discordActivity.AddActivity(activity)

		if err != nil {
			return err
		}
		time.Sleep(SLEEP_DELAY)
	}
	return nil
}

func (s *scrobbler) checkCurrentPlaying() *Listen {
	return s.listenBrainz.GetNowPlaying()
}

func NewScrobbler(username, clientId string) Scrobbler {
	return &scrobbler{
		username:        username,
		discordActivity: NewDiscordActivity(clientId),
		coverRetriever:  NewCoverArtRetriever(),
		listenBrainz:    NewListenBrainz(username),
	}
}
