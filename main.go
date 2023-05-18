package main

import (
	"flag"
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
	conf_location := flag.String("config-file", "conf.json", "Location of the config file")
	flag.Parse()
	config := LoadConfig(*conf_location)

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
