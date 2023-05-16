package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	scheduler := gocron.NewScheduler(time.UTC)
	var wg sync.WaitGroup

	job, err := scheduler.Every(15).Second().Do(func() {
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
