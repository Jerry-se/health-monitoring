package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"health-monitoring/db"
	hmp "health-monitoring/http"
	"health-monitoring/log"
	"health-monitoring/types"
	"health-monitoring/ws"

	"github.com/gin-gonic/gin"
)

var version string

func main() {
	configPath := flag.String("config", "", "run using the configuration file")
	versionFlag := flag.Bool("version", false, "show version number and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	if *configPath == "" {
		fmt.Println("run command like 'app -config ./config.json'")
		os.Exit(1)
	}
	cfg, err := types.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("Failed to load JSON configuration file:", err)
		os.Exit(1)
	}
	if err := log.InitLogrus(cfg.LogLevel, cfg.LogFile); err != nil {
		fmt.Println("Initialize the log failed:", err)
		os.Exit(1)
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := db.InitMongo(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database, cfg.MongoDB.ExpireTime); err != nil {
		os.Exit(1)
	}

	pm := hmp.NewPrometheusMetrics(cfg.Prometheus.JobName)

	router := gin.Default()
	router.GET("/metrics/prometheus", pm.Metrics)
	// router.GET("/echo", ws.Echo)
	router.GET("/websocket", func(c *gin.Context) {
		ws.Ws(c, pm)
	})

	// log.Log.Fatal(router.Run(cfg.Addr))

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Log.Fatalf("Start server: %v", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Log.Fatal("Server forced to shutdown: ", err)
	}

	log.Log.Println("Server exiting")
}
