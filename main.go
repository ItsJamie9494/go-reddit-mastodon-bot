package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func Fetch() *reddit.Post {
	client := reddit.DefaultClient()
	posts, _, err := client.Subreddit.HotPosts(context.Background(), "capybara", &reddit.ListOptions{
		Limit: 150, After: "", Before: "",
	})

	if err != nil {
		panic("handle this later lmao")
	}

	posts = Filter(posts, func(p *reddit.Post) bool {
		// TODO: Save image links so we don't have duplicates
		return !p.IsSelfPost && p.UpvoteRatio >= 0.9
	})

	var scores []int
	for _, s := range posts {
		scores = append(scores, s.Score)
	}
	posts = Filter(posts, func(p *reddit.Post) bool {
		best_score := Median(scores)
		return p.Score >= best_score
	})

	best_post := posts[rand.Intn(len(posts))]

	fmt.Printf(best_post.URL)

	return best_post
}

func main() {
	minute_schedule := flag.Int("schedule", 60, "How often to post (in minutes)")
	flag.Parse()

	scheduler := gocron.NewScheduler(time.UTC)
	var wg sync.WaitGroup

	job, err := scheduler.Every(*minute_schedule).Minutes().Do(func() {
		Fetch()
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
