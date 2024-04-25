package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

type model struct {
	progress     progress.Model
	spinner      spinner.Model
	listenBrainz ListenBrainz
	currentSong  TrackMetadata
	discord      DiscordActivity
	cover        CoverArtRetriever
}

type tickMsg time.Time
type trackChange TrackMetadata
type trackPlaying TrackMetadata
type trackStop struct{}

func (m model) getNowPlaying() *TrackMetadata {
	currentSong := m.listenBrainz.GetNowPlaying()

	if currentSong == nil {
		return nil
	}
	return &currentSong.TrackMetadata
}

func getCurrentSong(m model) tea.Msg {
	currentSong := m.getNowPlaying()
	if currentSong == nil {
		return trackStop{}
	}
	newSongRetriever := NewReleaseInfoRetriever(*currentSong)

	if !newSongRetriever.TrackMatches(m.currentSong) {
		return trackChange(*currentSong)
	}
	return trackPlaying(*currentSong)
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
		currentSong := TrackMetadata(msg)
		m.discord.Login()
		m.currentSong = currentSong
		log.Infof("Track changed to %s", m.getCurrentPlayingString())
		activity := m.getDiscordActivity()
		log.Debugf("activity: %v", activity)
		err := m.discord.AddActivity(activity)
		cobra.CheckErr(err)
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackPlaying:
		log.Debug("TRACK PLAYING")
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackStop:
		log.Debug("TRACK STOPPED")
		m.currentSong = TrackMetadata{}
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
	releaseInfo := NewReleaseInfoRetriever(m.currentSong)
	releaseId := releaseInfo.GetReleaseId()
	albumArt := m.cover.GetFront(releaseId)

	albumName := m.currentSong.ReleaseName
	artistName := m.currentSong.ArtistName
	trackName := m.currentSong.TrackName

	activity := &ScrobbleActivity{
		Album:     albumName,
		Artist:    artistName,
		Track:     trackName,
		Cover:     albumArt,
		ReleaseId: releaseId,
	}
	return activity
}

func (m model) getCurrentPlayingString() string {
	currentSong := m.currentSong
	trackName := currentSong.TrackName
	artistName := currentSong.ArtistName
	albumName := currentSong.ReleaseName
	return fmt.Sprintf("%s :: %s :: %s", trackName, artistName, albumName)
}

func (m model) View() string {
	currentSong := m.currentSong
	if currentSong.TrackName == "" {
		return fmt.Sprintf("%s No songs are playing", m.spinner.View())
	}
	return fmt.Sprintf("%s Playing %s", m.spinner.View(), m.getCurrentPlayingString())
}

func newSpinner() spinner.Model {
	return spinner.New(spinner.WithSpinner(spinner.Jump))
}

func NewModel() tea.Model {
	return &model{
		progress:     progress.New(progress.WithoutPercentage(), progress.WithDefaultGradient()),
		spinner:      newSpinner(),
		listenBrainz: NewListenBrainz(),
		discord:      NewDiscordActivity(),
		cover:        NewCoverArtRetriever(),
	}
}
