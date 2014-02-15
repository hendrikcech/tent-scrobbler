package track

import (
	"time"
)

type Track struct {
	ID        string
	Name      string
	Album     string
	Artist    string
	Remaining time.Duration
	Duration  time.Duration
}
