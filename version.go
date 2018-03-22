package nasher

import (
	"time"
)

var (
	Version   string
	Commit    string
	BuildTime string
	Title     string
	StartTime time.Time
)

func init() {
	if Version == "" {
		Version = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if Title == "" {
		Title = "nasher"
	}
	StartTime = time.Now()
}
