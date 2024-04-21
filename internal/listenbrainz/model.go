package listenbrainz

type ListenResponse struct {
	Payload ListenPayload `json:"payload"`
}

type ListenPayload struct {
	Count      int      `json:"count"`
	PlayingNow bool     `json:"playing_now"`
	UserId     string   `json:"user_id"`
	Listens    []Listen `json:"listens"`
}

type Listen struct {
	PlayingNow    bool          `json:"playing_now"`
	ListenedAt    int           `json:"listened_at"`
	TrackMetadata TrackMetadata `json:"track_metadata"`
}

type TrackMetadata struct {
	AdditionalInfo TrackAdditionalInfo `json:"additional_info"`
	ArtistName     string              `json:"artist_name"`
	TrackName      string              `json:"track_name"`
	ReleaseName    string              `json:"release_name"`
	MbidMapping    MbidMapping         `json:"mbid_mapping"`
}

type TrackAdditionalInfo struct {
	MbidMapping
	TrackNumber int      `json:"track_number"`
	Tags        []string `json:"tags"`
	Duration    int      `json:"duration"`
	DurationMs  int      `json:"duration_ms"`
}

type MbidMapping struct {
	ReleaseMBID      string   `json:"release_mbid"`
	ArtistMBIDs      []string `json:"artist_mbids"`
	ReleaseGroupMBID string   `json:"release_group_mbid"`
	RecordingMBID    string   `json:"recording_mbid"`
	TrackMBID        string   `json:"track_mbid"`
}
