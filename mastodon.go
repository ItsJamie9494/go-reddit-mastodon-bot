package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mattn/go-mastodon"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type MediaAttachments struct {
	ImageSizeLimit     int      `json:"image_size_limit"`
	ImageMatrixLimit   int      `json:"image_matrix_limit"`
	SupportedMimeTypes []string `json:"supported_mime_types"`
}

type Configuration struct {
	Attachments MediaAttachments `json:"media_attachments"`
}

type Config struct {
	Config Configuration `json:"configuration"`
}

func GetMastodonClientWithLimits(GlobalConfig *GlobalConfig) (*mastodon.Client, Config) {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       GlobalConfig.APIBaseURL,
		ClientID:     GlobalConfig.ClientID,
		ClientSecret: GlobalConfig.ClientSecret,
		AccessToken:  GlobalConfig.AccessToken,
	})

	var data *Config
	url, _ := url.Parse(c.Config.Server)
	url.Path = "api/v2/instance"
	req, err := http.Get(url.String())
	if err != nil {
		log.Fatal("Unable to make request to Mastodon server, please check that the URL is correct", ": ", err)
	}
	body, _ := io.ReadAll(req.Body)
	req.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Unable to parse server response on endpoint /api/v2/instance", ": ", err)
	}

	return c, *data
}

func Fetch(GlobalConfig *GlobalConfig, ImgLocation string) *reddit.Post {
	client := reddit.DefaultClient()
	images := LoadImagesFile(ImgLocation)
	posts, _, err := client.Subreddit.HotPosts(context.Background(), GlobalConfig.Subreddit, &reddit.ListOptions{
		Limit: 150, After: "", Before: "",
	})

	if err != nil {
		panic("handle this later lmao")
	}

	posts = Filter(posts, func(p *reddit.Post) bool {
		return !p.IsSelfPost && p.UpvoteRatio >= 0.9 && !Contains(images, p.URL)
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

	log.Default().Print(best_post.URL)
	AppendToImagesFile(ImgLocation, best_post.URL)

	return best_post
}

func ValidateMedia(c Config, URL string) bool {
	req, err := http.Head(URL)
	if err != nil {
		log.Fatal("Unable to make HEAD request to endpoint ", URL, ": ", err)
	}
	str, _ := strconv.Atoi(req.Header.Get("Content-Length"))
	if str <= c.Config.Attachments.ImageSizeLimit && Contains(c.Config.Attachments.SupportedMimeTypes, req.Header.Get("Content-Type")) {
		return true
	} else {
		return false
	}
}

func UploadMedia(GlobalConfig *GlobalConfig, ImgLocation string) {
	client, config := GetMastodonClientWithLimits(GlobalConfig)
	image := Fetch(GlobalConfig, ImgLocation)
	for !ValidateMedia(config, image.URL) {
		image = Fetch(GlobalConfig, ImgLocation)
	}
	req, err := http.Get(image.URL)
	if err != nil {
		log.Fatal("Unable to make GET request to endpoint ", image.URL, ": ", err)
	}
	body, _ := io.ReadAll(req.Body)
	req.Body.Close()
	attachment, err := client.UploadMediaFromBytes(context.Background(), body)
	if err != nil {
		log.Fatal("Unable to upload media to Mastodon: ", err)
	}
	status, err := client.PostStatus(context.Background(), &mastodon.Toot{
		Status:      fmt.Sprintf("%s (by u/%s)", image.Title, image.Author),
		InReplyToID: "",
		MediaIDs:    []mastodon.ID{attachment.ID},
		Sensitive:   false,
		SpoilerText: "",
		Visibility:  "",
		Language:    "",
	})
	if err != nil {
		log.Fatal("Unable to post status to Mastodon: ", err)
	}
	log.Default().Print(status.URL)
}
