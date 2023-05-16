package main

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	config := LoadConfig()
	scheduler := gocron.NewScheduler(time.UTC)
	var wg sync.WaitGroup

	job, err := scheduler.Every(config.Schedule).Minutes().Do(func() {
		UploadMedia(config)
	})
	if err != nil {
		panic("Unable to start scheduler")
	}

	scheduler.StartAsync()

	for !job.IsRunning() {
		wg.Add(1)
		wg.Wait()
	}
}
