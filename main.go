package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
)

func MainLoop() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func main() {
	config := LoadConfig()

	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(config.Schedule).Minutes().Do(func() {
		if (config.Schedule == 60 && time.Now().Minute() == 0) || config.Schedule != 60 {
			UploadMedia(config)
		}
	})
	if err != nil {
		panic("Unable to start scheduler")
	}

	scheduler.StartAsync()

	MainLoop()
}
