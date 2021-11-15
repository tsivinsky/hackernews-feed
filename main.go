package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/0xAX/notificator"
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

var (
	allowedCommands = []string{"new", "best", "top", "ask"}
	command         = "new"
)

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		command = args[0]
	}

	if !isCommandAllowed(command) {
		fmt.Printf("Command %s doesn't exist\n", command)
		os.Exit(1)
	}

	storyChan := make(chan Story)
	cmdChan := make(chan string)

	go listenForStories(apiUrl, storyChan)
	go listenForCommands(cmdChan)

	var story Story
	for {
		select {
		case story = <-storyChan:
			printStory(story)
		case cmd := <-cmdChan:
			if cmd == "open" || cmd == "o" {
				err := openInBrowser(story.Url)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func listenForCommands(c chan string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		s = strings.TrimSpace(s)

		c <- s
	}
}

func listenForStories(apiUrl string, c chan Story) {
	notify := notificator.New(notificator.Options{
		AppName: "hacker-news-feed",
	})

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

		if latestStoryId > 0 {
			notify.Push("New Hacker Story", story.Title, "", notificator.UR_NORMAL)
		}

		if story.Url == "" {
			story.Url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.Id)
			story.Title = fmt.Sprintf("[Ask] %s", story.Title)
		}

		latestStoryId = story.Id

		c <- story
	}
}

func openInBrowser(url string) error {
	cmd := exec.Command("xdg-open", url)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func printStory(story Story) {
	var storyInString string

	storyInString += "\""
	storyInString += string(YELLOW_COLOR)
	storyInString += fmt.Sprintf("%s", story.Title)
	storyInString += string(RESET_COLOR)
	storyInString += "\""
	storyInString += fmt.Sprintf(" by %s\n", story.By)
	storyInString += fmt.Sprintf("  - %s\n", story.Url)

	fmt.Print(storyInString)
}

func isCommandAllowed(command string) bool {
	for _, cmd := range allowedCommands {
		if cmd == command {
			return true
		}
	}

	return false
}
