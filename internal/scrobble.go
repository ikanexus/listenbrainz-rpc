package internal

import (
	"time"

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
	CheckCurrentPlaying() *Listen
}

func (s scrobbler) Scrobble() error {
	logger.Infof("Running with user: %s", s.username)
	defer s.discordActivity.Logout()
	s.discordActivity.Login()

	isIdle := false

	var currentSongId string
	currentSongStart := time.Now()
	var runtime int64

	for running := true; running; {
		currentSong := s.CheckCurrentPlaying()

		if currentSong == nil {
			time.Sleep(SLEEP_DELAY * 2)
			if !isIdle {
				logger.Infof("No song playing - clearing activity")
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
			logger.Debugf("Still playing %s :: %s :: %s", trackName, artistName, albumName)
			startDuration := (time.Duration(runtime) * time.Millisecond) + SLEEP_DELAY
			duration := currentSongStart.Add(startDuration).Sub(currentSongStart)
			logger.Debugf("Estimated runtime: %v/%v", duration, trackDuration)
			runtime = startDuration.Milliseconds()
			time.Sleep(SLEEP_DELAY)
			continue
		}
		currentSongStart = time.Now()
		runtime = 0
		logger.Infof("Playing %s :: %s :: %s", trackName, artistName, albumName)
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

func (s *scrobbler) CheckCurrentPlaying() *Listen {
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
