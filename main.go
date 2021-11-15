package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/tsivinsky/hackernews-feed/api"
)

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

var (
	intervalTime int = 2
)

func main() {
	args := os.Args[1:]

	if i, exists := isItemExistInSlice(args, "-i"); exists {
		intervalTime, _ = strconv.Atoi(args[i+1])
		args = append(args[:i], args[i+2:]...) // Go Team, please add method for removing elements from slice, pls, pls, pls
	}

	if len(args) > 0 {
		command = args[0]
	}

	if _, exists := isItemExistInSlice(allowedCommands, command); !exists {
		fmt.Printf("Command %s doesn't exist\n", command)
		os.Exit(1)
	}

	apiUrl := fmt.Sprintf("%s/%sstories.json", api.BaseApiUrl, command)

	storyChan := make(chan api.Story)
	cmdChan := make(chan string)

	go listenForStories(apiUrl, storyChan)
	go listenForCommands(cmdChan)

	var story api.Story
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

func listenForCommands(cmdChan chan string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		s = strings.TrimSpace(s)

		cmdChan <- s
	}
}

func listenForStories(apiUrl string, storyChan chan api.Story) {
	notify := notificator.New(notificator.Options{
		AppName: "hacker-news-feed",
	})

	var latestStoryId int

	for {
		newStories, err := api.GetNewStories(apiUrl)
		if err != nil {
			log.Fatal(err)
		}
		newStoryId := newStories[0]

		if newStoryId == latestStoryId {
			log.Println("Going to sleep")
			time.Sleep(time.Duration(intervalTime) * time.Second)
			continue
		}

		story, err := api.GetStoryById(newStoryId)
		if err != nil {
			log.Fatal(err)
		}

		// Sometimes api returns empty story, and at this point I'm too afraid to ask
		if story.Id == 0 {
			continue
		}

		if latestStoryId > 0 {
			notify.Push("New Hacker Story", story.Title, "", notificator.UR_NORMAL)
		}

		// [Ask] stories don't have external url
		if story.Url == "" {
			story.Url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.Id)
			story.Title = fmt.Sprintf("[Ask] %s", story.Title)
		}

		latestStoryId = story.Id

		storyChan <- story
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

func printStory(story api.Story) {
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

func isItemExistInSlice(slice interface{}, item interface{}) (int, bool) {
	s := reflect.ValueOf(slice)

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return i, true
		}
	}

	return -1, false
}
