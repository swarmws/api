package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

var ServerStartTime time.Time

func SetupRoutes(r *gin.Engine) {
	routes(r)
}

func routes(r *gin.Engine) {
	Regitserdomains()
	// public health + research-facing route groups only —
	// full production route tree not included in this showcase
	r.GET("/c", health)
	_ = r.Group("/api")
}

func health(c *gin.Context) {
	c.String(200, "ok")
}

func HealthCheck(c *gin.Context) {
	health(c)
}

func CheckExternalCovers() (ext, noValid int) { return 0, 0 }
func FixCovers()                             {}
