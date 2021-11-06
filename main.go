package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Story struct {
	Id    int    `json:"id"'`
	Title string `json:"title"`
	Url   string `json:"url"`
	By    string `json:"by"`
}

func main() {
	var latestStoryId int

	for {
		newStories, err := getNewStories()
		if err != nil {
			log.Fatal(err)
		}
		newStoryId := newStories[0]

		if newStoryId == latestStoryId {
			time.Sleep(2 * time.Second)
			continue
		}

		story, err := getStoryById(newStoryId)
		if err != nil {
			log.Fatal(err)
		}

		latestStoryId = story.Id

		fmt.Printf("\"%s\" by %s\n", story.Title, story.By)

		var url string
		if story.Url != "" {
			url = story.Url
		} else {
			url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.Id)
		}

		fmt.Printf("  %s\n", url)
	}
}

func getNewStories() ([]int, error) {
	var items []int

	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/newstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func getStoryById(storyId int) (Story, error) {
	var story Story

	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyId)
	resp, err := http.Get(url)
	if err != nil {
		return Story{}, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&story)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}
