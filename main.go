package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"swarm-api/api"
	"swarm-api/api/data"
	"swarm-api/api/db"
	"swarm-api/api/noah"
	"swarm-api/api/scrapes"
	"swarm-api/api/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	api.ServerStartTime = time.Now()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) { c.Next() })
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PUT", "PATCHd"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "User-Agent", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	db.InitializeDB()
	db.EnsureIndexes()
	db.InitializeRedis()
	db.InitializeWasabi()
	scrapes.NoahSaveChapter = noah.SaveChapterToNoah
	scrapes.NoahUploadChapterPanels = noah.UploadChapterImageBytes
	data.StartSimpleSync()
	db.OnChapterInserted = scrapes.UploadNewChapterPanelsToNoahIfDownload
	db.StartChapterInsertSummary(5 * time.Minute)
	paused := utils.UpdatesPollerDisabled()
	if !paused {
		utils.EnsureCFBypassServer()
	}
	disablePollers := "false"
	if paused {
		disablePollers = "true"
	}
	pollerStarted, pollerTotal := api.StartPollers(disablePollers)
	go scrapes.StartFinished()
	go func() {
		time.Sleep(1 * time.Minute)
		if err := db.RefreshMangaBookmarkCounts(); err != nil {
			log.Printf("Error refreshing bb counts: %v", err)
		}
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := db.RefreshMangaBookmarkCounts(); err != nil {
				log.Printf("Error refreshing bb counts: %v", err)
			}
		}
	}()
	api.SetupRoutes(r)
	go func() {
		time.Sleep(15 * time.Second)
		ext, noValid := api.CheckExternalCovers()
		log.Printf("covers: %d external, %d without valid cum.swarm.ws cover", ext, noValid)
		api.FixCovers()
	}()
	go noah.StartNoahServer()
	if !paused {
		// noah.SaveAll()
	}
	port := utils.HTTPListenPort()
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   10 * time.Minute,
		IdleTimeout:    2 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("\033[48;2;245;0;255m %s \033[0m swarm-api up · http :%s · noah :%s · poller :%d/%d\n",
		time.Now().Format("2006/01/02 15:04:05"), port, noah.NoahProxyPort, pollerStarted, pollerTotal)
	{
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		api.HealthCheck(c)
		if w.Code != http.StatusOK {
			log.Printf("startup health-check: degraded %s", strings.TrimSpace(w.Body.String()))
		}
	}
	db.StartMangaTableMirrorBootstrap()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
