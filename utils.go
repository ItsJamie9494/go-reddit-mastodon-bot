package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type GlobalConfig struct {
	ClientID     string
	ClientSecret string
	APIBaseURL   string
	Subreddit    string
	Schedule     int
}

func LoadConfig() *GlobalConfig {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := GlobalConfig{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatal("Failed to load configuration file: ", err)
	}
	return &conf
}

func Median(data []int) int {
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)

	sort.Ints(dataCopy)

	var median int
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Contains[T comparable](ss []T, elem T) bool {
	for _, v := range ss {
		if v == elem {
			return true
		}
	}

	return false
}
