package internal

import (
	"encoding/json"
)

type ReleaseInfoRetriever interface {
	GetDuration() int
	GetReleaseId() string
	GetTrackHash() string
}

type releaseInfoRetriever struct {
	track TrackMetadata
}

func NewReleaseInfoRetriever(track TrackMetadata) ReleaseInfoRetriever {
	return &releaseInfoRetriever{
		track: track,
	}
}

func (r releaseInfoRetriever) GetDuration() int {
	trackInfo := r.track.AdditionalInfo
	if trackInfo.Duration != 0 {
		return trackInfo.Duration
	}
	if trackInfo.DurationMs != 0 {
		return trackInfo.DurationMs / 1000
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

func (r releaseInfoRetriever) GetTrackHash() string {
	b, _ := json.Marshal(r.track)
	return string(b)
}
