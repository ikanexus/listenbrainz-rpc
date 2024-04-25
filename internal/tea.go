package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	progress    progress.Model
	spinner     spinner.Model
	scrobbler   Scrobbler
	currentSong TrackMetadata
}

type tickMsg time.Time
type trackChange TrackMetadata
type trackStop struct{}

func (m model) getNowPlaying() *TrackMetadata {
	currentSong := m.scrobbler.CheckCurrentPlaying()
	if currentSong == nil {
		return nil
	}
	return &currentSong.TrackMetadata
}

func tickCmd(m model) tea.Cmd {
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		currentSong := m.getNowPlaying()
		if currentSong == nil {
			return trackStop{}
		}
		newSongRetriever := NewReleaseInfoRetriever(*currentSong)
		if !newSongRetriever.TrackMatches(m.currentSong) {
			return trackChange(*currentSong)
		}
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	log.Printf("test")
	return tickCmd(m)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		// if m.progress.Percent() == 1.0 {
		// 	return m, tea.Quit
		// }

		// cmd := m.progress.IncrPercent(0.1)
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackChange:
		currentSong := TrackMetadata(msg.(trackChange))
		m.currentSong = currentSong
		return m, tea.Batch(tickCmd(m), m.spinner.Tick)
	case trackStop:
		m.currentSong = TrackMetadata{}
		m.spinner = newSpinner()
		return m, tea.Batch(tickCmd(m))
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	currentSong := m.currentSong
	if currentSong.TrackName == "" {
		return fmt.Sprintf("%s No songs are playing", m.spinner.View())
	}
	return fmt.Sprintf("%s Playing %s :: %s :: %s", m.spinner.View(), currentSong.TrackName, currentSong.ArtistName, currentSong.ReleaseName)
	// return m.progress.View()
}

func newSpinner() spinner.Model {
	return spinner.New(spinner.WithSpinner(spinner.Jump))
}

func NewModel() tea.Model {
	return &model{
		progress:  progress.New(progress.WithoutPercentage(), progress.WithDefaultGradient()),
		spinner:   newSpinner(),
		scrobbler: NewScrobbler(),
	}
}
