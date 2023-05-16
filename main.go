package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	minute_schedule := flag.Int("schedule", 60, "How often to post (in minutes)")
	flag.Parse()

	scheduler := gocron.NewScheduler(time.UTC)
	var wg sync.WaitGroup

	job, err := scheduler.Every(*minute_schedule).Minutes().Do(func() {
		fmt.Println("Meow")
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
