package db

import "time"

type Manga struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	Cover     string   `json:"cover,omitempty"`
	Covers    []string `json:"covers"`
	Status    string   `json:"status,omitempty"`
	Type      string   `json:"type,omitempty"`
}

type Chapter struct {
	ID      string    `json:"id"`
	MangaID string    `json:"manga_id,omitempty"`
	Number  string    `json:"number"`
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Source  string    `json:"source,omitempty"`
	Pages   []string  `json:"pages"`
	Hash    string    `json:"hash,omitempty"`
	Date    time.Time `json:"date,omitempty"`
}

func InitializeDB()                                      {}
func EnsureIndexes()                                     {}
func StartChapterInsertSummary(interval time.Duration)   {}
func StartMangaTableMirrorBootstrap()                    {}
func RefreshMangaBookmarkCounts() error                  { return nil }

var OnChapterInserted func(mangaID, chapterID string)
