package scrapes

import (
	"time"

	"swarm-api/api/data/pollclock"
)

type SourcePanelsScraper func(string) ([]string, string, error)

type SourcePoller struct {
	Function    string `json:"function"`
	SiteURL     string `json:"site_url"`
	PollTimeMin int    `json:"poll_time_min"`
	SourceName  string `json:"source_name"`
	Priority    int    `json:"priority"`
	UpdateURL   string `json:"update_url"`
	MangaURL    string `json:"manga_url"`
	Chap        string `json:"chap_url"`
	CoverURL    string `json:"cover_url"`
	Panels      string `json:"chapter_panels_urls_from_html"`
	Proxy       bool   `json:"proxy"`
	CFServer    bool   `json:"cf_server"`
	Cookies     string `json:"cookies"`
}

func GetSources() []SourcePoller {
	return nil // configured sources live in production; omitted here
}

func (s SourcePoller) PollInterval() time.Duration {
	if s.PollTimeMin <= 0 {
		return pollclock.Cap(time.Minute)
	}
	return pollclock.Cap(time.Duration(s.PollTimeMin) * time.Minute)
}

func DefaultPanels(u string) ([]string, string, error) { return nil, "", nil }
func ConfiguredPollerSourceTally() (started, total int) { return 0, 0 }
func StartPoller() (started, total int)                 { return 0, 0 }
func StartDebugDiscordDailyDigest()                     {}

// site panel / Start* symbols referenced by source.go (stubs)
func AsuraPanels(string) ([]string, string, error)  { return nil, "", nil }
func MangaDexPanels(string) ([]string, string, error) { return nil, "", nil }
func FlamePanels(string) ([]string, string, error)  { return nil, "", nil }
func StartAsura()                                   {}
func StartMangaDex()                                {}
func StartFlameScans()                              {}

