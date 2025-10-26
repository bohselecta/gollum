package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/haydenlabs/gollum/apis/openai"
	"github.com/haydenlabs/gollum/engine/impl"
	"github.com/haydenlabs/gollum/metrics"
	"github.com/haydenlabs/gollum/obs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Register metrics
	metrics.MustRegister()
	
	engine := impl.NewEngine()
	r := gin.Default()
	obs.Install(r)

	// Static assets and favicon
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/icons/gollum_icon_256.png")
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := openai.NewAPI(engine)
	api.Register(r)

	// selfTestMetal() // Disabled - Metal files temporarily disabled

	// Minimal landing page
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, landingHTML)
	})
	log.Printf("GoLLuM listening on :%s", port)
	go func() { log.Println(http.ListenAndServe(":6060", nil)) }()
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
