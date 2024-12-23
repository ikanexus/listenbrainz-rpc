package internal

import (
	"encoding/json"
	"time"
)

type ReleaseInfoRetriever interface {
	GetDuration() time.Duration
	GetReleaseId() string
	GetTrackId() string
	GetTrackHash() string
	TrackMatches(other TrackMetadata) bool
}

type releaseInfoRetriever struct {
	track TrackMetadata
}

func NewReleaseInfoRetriever(track TrackMetadata) ReleaseInfoRetriever {
	return &releaseInfoRetriever{
		track: track,
	}
}

func (r releaseInfoRetriever) GetDuration() time.Duration {
	trackInfo := r.track.AdditionalInfo
	if trackInfo.Duration != 0 {
		return time.Duration(trackInfo.Duration) * time.Second
	}
	if trackInfo.DurationMs != 0 {
		return time.Duration(trackInfo.DurationMs) * time.Millisecond
	}
	return 0
}

func (r releaseInfoRetriever) GetReleaseId() string {
	metadata := r.track
	locations := []MbidMapping{
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

func (r releaseInfoRetriever) GetTrackId() string {
	metadata := r.track
	locations := []MbidMapping{
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

func (r releaseInfoRetriever) GetTrackHash() string {
	b, _ := json.Marshal(r.track)
	return string(b)
}

func (r releaseInfoRetriever) TrackMatches(b TrackMetadata) bool {
	self := r.GetTrackHash()
	other := NewReleaseInfoRetriever(b).GetTrackHash()

	return self == other
}
