package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"health-monitoring/db"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.InitMongo(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database, cfg.MongoDB.ExpireTime); err != nil {
		os.Exit(1)
	}

	router := gin.Default()
	router.GET("/echo", ws.Echo)
	router.GET("/websocket", ws.Ws)
	log.Log.Fatal(router.Run(cfg.Addr))
}
