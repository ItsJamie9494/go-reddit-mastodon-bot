package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type GlobalConfig struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	APIBaseURL   string
	Subreddit    string
	Schedule     int
}

func LoadConfig(filename string) *GlobalConfig {
	file, _ := os.Open(filename)
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

func AppendToImagesFile(File string, URL string) {
	file, err := os.OpenFile(File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Default().Print("WARNING: Cannot access images.txt. Does the file exist?")
		return
	}
	defer file.Close()

	if _, err = file.WriteString(fmt.Sprint(URL, "\n")); err != nil {
		log.Default().Print("WARNING: Failed to save current image to ", File)
		return
	}
}

func LoadImagesFile(File string) []string {
	file, err := os.OpenFile(File, os.O_RDONLY, 0600)
	if err != nil {
		log.Default().Print("WARNING: Cannot access images.txt. Does the file exist?")
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		log.Default().Print("WARNING: Failed to read ", File)
		return nil
	}
	return lines
}
