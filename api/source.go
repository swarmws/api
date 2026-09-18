package api

import (
	"net/url"
	"strings"

	"swarm-api/api/scrapes"
)

// SourceDomainMap maps scrape hosts to panel extractors.
var SourceDomainMap = map[string]any{
	"asurascans.com":  scrapes.AsuraPanels,
	"mangadex.org":    scrapes.MangaDexPanels,
	"flamecomics.xyz": scrapes.FlamePanels,
	// … additional registered sources omitted in this showcase slice
}

var siteUpdatePollerFuncs = []func(){
	scrapes.StartAsura,
	scrapes.StartMangaDex,
	scrapes.StartFlameScans,
	// …
}

func Regitserdomains() {
	for _, source := range scrapes.GetSources() {
		fallback := scrapes.DefaultPanels
		for _, raw := range []string{source.SiteURL, source.UpdateURL, source.MangaURL} {
			u, err := url.Parse(strings.TrimSpace(raw))
			if err != nil || u.Host == "" {
				continue
			}
			host := strings.ToLower(strings.TrimPrefix(u.Host, "www."))
			if _, exists := SourceDomainMap[host]; exists {
				continue
			}
			SourceDomainMap[host] = fallback
		}
	}
}

func StartPollers(disableUpdates string) (pollerStarted, pollerTotal int) {
	_, cfgTotal := scrapes.ConfiguredPollerSourceTally()
	if strings.EqualFold(strings.TrimSpace(disableUpdates), "true") {
		return 0, cfgTotal + 1 + len(siteUpdatePollerFuncs)
	}

	pollerTotal++
	scrapes.StartDebugDiscordDailyDigest()
	pollerStarted++

	cs, ct := scrapes.StartPoller()
	pollerStarted += cs
	pollerTotal += ct

	for _, fn := range siteUpdatePollerFuncs {
		pollerTotal++
		fn()
		pollerStarted++
	}
	return pollerStarted, pollerTotal
}
