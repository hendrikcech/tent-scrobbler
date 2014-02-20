package example

import (
	"github.com/hendrikcech/tent-scrobbler/track"
	"time"
)

func CurrentTrack() (song track.Track, err error) {
	song = track.Track{
		ID:        "ID",
		Name:      "Name",
		Album:     "Album",
		Artist:    "Artist",
		Remaining: time.Duration(1200) * time.Second,
		Duration:  time.Duration(1800) * time.Second,
	}

	return
}
