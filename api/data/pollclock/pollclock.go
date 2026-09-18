package pollclock

import "time"

func Cap(d time.Duration) time.Duration {
	return d
}

func RunLoopFromStart(interval time.Duration, run func()) {}
