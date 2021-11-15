package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	Id    int    `json:"id"'`
	Title string `json:"title"`
	Url   string `json:"url"`
	By    string `json:"by"`
}

var BaseApiUrl = "https://hacker-news.firebaseio.com/v0"

func GetNewStories(apiUrl string) ([]int, error) {
	var items []int

	resp, err := http.Get(apiUrl)
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

// GetStoryById takes storyId as an argument and returns story or error
func GetStoryById(storyId int) (Story, error) {
	var story Story

	url := fmt.Sprintf("%s/item/%d.json", BaseApiUrl, storyId)
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
