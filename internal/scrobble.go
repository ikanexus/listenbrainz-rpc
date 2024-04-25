package internal

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
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
	currentSongStart := time.Now()
	var runtime int64

	for running := true; running; {
		currentSong := s.checkCurrentPlaying()

		if currentSong == nil {
			time.Sleep(SLEEP_DELAY * 2)
			if !isIdle {
				log.Infof("No song playing - clearing activity")
				isIdle = true
				s.discordActivity.Logout()
			}
			continue
		}

		if isIdle {
			s.discordActivity.Login()
			isIdle = false
		}

		releaseInfo := NewReleaseInfoRetriever(currentSong.TrackMetadata)

		trackDuration := releaseInfo.GetDuration()
		trackHash := releaseInfo.GetTrackHash()
		trackMetadata := currentSong.TrackMetadata
		trackName := trackMetadata.TrackName
		artistName := trackMetadata.ArtistName
		albumName := trackMetadata.ReleaseName

		if trackHash == currentSongId {
			log.Debugf("Still playing %s :: %s :: %s", trackName, artistName, albumName)
			startDuration := (time.Duration(runtime) * time.Millisecond) + SLEEP_DELAY
			duration := currentSongStart.Add(startDuration).Sub(currentSongStart)
			log.Debugf("Estimated runtime: %v/%v", duration, trackDuration)
			runtime = startDuration.Milliseconds()
			time.Sleep(SLEEP_DELAY)
			continue
		}
		currentSongStart = time.Now()
		runtime = 0
		log.Infof("Playing %s :: %s :: %s", trackName, artistName, albumName)
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

func NewScrobbler() Scrobbler {
	username := viper.GetString("user")
	return &scrobbler{
		username:        username,
		discordActivity: NewDiscordActivity(),
		coverRetriever:  NewCoverArtRetriever(),
		listenBrainz:    NewListenBrainz(username),
	}
}
