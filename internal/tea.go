package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
	progress     progress.Model
	spinner      spinner.Model
	listenBrainz ListenBrainz
	currentListen Listen
	discord      DiscordActivity
	cover        CoverArtRetriever
}

type (
	trackChange  Listen
	trackPlaying Listen
	trackStop    struct{}
)

func (m model) getNowPlaying() *Listen {
	return m.listenBrainz.GetNowPlaying()
}

func getCurrentSong(m model) tea.Msg {
	currentListen := m.getNowPlaying()
	if currentListen == nil {
		return trackStop{}
	}
	newSongRetriever := NewReleaseInfoRetriever(currentListen.TrackMetadata)

	if !newSongRetriever.TrackMatches(m.currentListen.TrackMetadata) {
		return trackChange(*currentListen)
	}
	return trackPlaying(*currentListen)
}

func tickCmd(m model) tea.Cmd {
	return tea.Tick(time.Second*10, func(_ time.Time) tea.Msg {
		return getCurrentSong(m)
	})
}

func (m model) Init() tea.Cmd {
	currentSong := func() tea.Msg {
		return getCurrentSong(m)
	}
	defer m.discord.Logout()
	return currentSong
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case trackChange:
		log.Debug("TRACK CHANGE")
		currentListen := Listen(msg)
		m.discord.Login()
		m.currentListen = currentListen
		log.Infof("Track changed to %s", m.getCurrentPlayingString())
		activity := m.getDiscordActivity()
		log.Debugf("activity: %v", activity)
		err := m.discord.AddActivity(activity)
		if err != nil {
			log.Errorf("Unable to add Discord activity: %v", err)
		}
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackPlaying:
		log.Debug("TRACK PLAYING")
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackStop:
		log.Debug("TRACK STOPPED")
		m.currentListen = Listen{}
		m.spinner = newSpinner()
		m.discord.Logout()
		return m, tickCmd(m)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) getDiscordActivity() *ScrobbleActivity {
	releaseInfo := NewReleaseInfoRetriever(m.currentListen.TrackMetadata)
	releaseId := releaseInfo.GetReleaseId()
	albumArt := m.cover.GetFront(releaseId)

	albumName := m.currentListen.TrackMetadata.ReleaseName
	artistName := m.currentListen.TrackMetadata.ArtistName
	trackName := m.currentListen.TrackMetadata.TrackName
	duration := releaseInfo.GetDuration()

	activity := &ScrobbleActivity{
		Album:     albumName,
		Artist:    artistName,
		Track:     trackName,
		Cover:     albumArt,
		ReleaseId: releaseId,
		Duration:  duration,
	}
	return activity
}

func (m model) getCurrentPlayingString() string {
	trackName := m.currentListen.TrackMetadata.TrackName
	artistName := m.currentListen.TrackMetadata.ArtistName
	albumName := m.currentListen.TrackMetadata.ReleaseName
	return fmt.Sprintf("%s :: %s :: %s", trackName, artistName, albumName)
}

func (m model) View() string {
	if m.currentListen.TrackMetadata.TrackName == "" {
		return fmt.Sprintf("%s No songs are playing", m.spinner.View())
	}
	return fmt.Sprintf("%s Playing %s", m.spinner.View(), m.getCurrentPlayingString())
}

func newSpinner() spinner.Model {
	return spinner.New(spinner.WithSpinner(spinner.Jump))
}

func NewModel(cfg Config) tea.Model {
	return &model{
		progress:     progress.New(progress.WithoutPercentage(), progress.WithDefaultGradient()),
		spinner:      newSpinner(),
		listenBrainz: NewListenBrainz(cfg),
		discord:      NewDiscordActivity(cfg),
		cover:        NewCoverArtRetriever(),
	}
}
