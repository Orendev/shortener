package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/Orendev/shortener/internal/app"
	"github.com/Orendev/shortener/internal/config"
	"github.com/Orendev/shortener/internal/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Build go build -ldflags "-X main.buildVersion=1.0.1 -X main.buildCommit=1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')'" -o shortener  cmd/shortener/main.go
// Run: ./shortener
func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.NewLogger(cfg.Log.FlagLogLevel); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	app.Run(cfg)
}
