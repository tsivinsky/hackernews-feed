package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Story struct {
	Id    int    `json:"id"'`
	Title string `json:"title"`
	Url   string `json:"url"`
	By    string `json:"by"`
}

const (
	RESET_COLOR  = "\033[0m"
	RED_COLOR    = "\033[31m"
	GREEN_COLOR  = "\033[32m"
	YELLOW_COLOR = "\033[33m"
	BLUE_COLOR   = "\033[34m"
	PURPLE_COLOR = "\033[35m"
	CYAN_COLOR   = "\033[36m"
	WHITE_COLOR  = "\033[37m"
)

var command = "new"

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		command = args[0]
	}

	var apiUrl = "https://hacker-news.firebaseio.com/v0"

	switch command {
	case "new":
		apiUrl = apiUrl + "/newstories.json"
		break
	case "best":
		apiUrl = apiUrl + "/beststories.json"
		break
	case "top":
		apiUrl = apiUrl + "/topstories.json"
		break
	case "ask":
		apiUrl = apiUrl + "/askstories.json"
		break
	default:
		fmt.Printf("Command %s does not exist\n", command)
		os.Exit(1)
	}

	var latestStoryId int

	for {
		newStories, err := getNewStories(apiUrl)
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

		if story.Id == 0 {
			continue
		}

		latestStoryId = story.Id

		fmt.Printf("\"")
		fmt.Print(string(YELLOW_COLOR))
		fmt.Printf("%s", story.Title)
		fmt.Print(string(RESET_COLOR))
		fmt.Print("\"")
		fmt.Printf(" by %s\n", story.By)

		var url string
		if story.Url != "" {
			url = story.Url
		} else {
			url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.Id)
		}

		fmt.Printf("  - %s\n", url)
	}
}

func getNewStories(apiUrl string) ([]int, error) {
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
