package scrapes

func StartFinished() {}

func UploadNewChapterPanelsToNoahIfDownload(mangaID, chapterID string) {}

var NoahSaveChapter func(chapterURL, mangaID, chapterNumber string, imageURLs []string) ([]string, error)
var NoahUploadChapterPanels func(chapterID, mangaID string, raw [][]byte) ([]string, error)
