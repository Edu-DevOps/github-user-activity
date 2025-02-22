package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"io"
	"flag"

	"github-user-activity/models"
)

var (
	user string
)

func main() {
	flag.StringVar(&user, "user", "", "Github username")
	flag.Parse()
	if user == "" {
		fmt.Println("--user flag is required")
		os.Exit(1)
	}

	gitToken := os.Getenv("GITHUB_TOKEN")
	events, _ := getEvents(user, gitToken)
	translatedEvents, _ := translateEvents(events)
	outputInit := "Output:"
	fmt.Println(outputInit)
	for index := range translatedEvents {
		outputBody := "  - " + translatedEvents[index].Type + ": "
		if translatedEvents[index].Type == "Suggestion Applied" {
			outputBody := outputBody + "\n    URL: " + translatedEvents[index].Payload.PullRequest.URL
			fmt.Println(outputBody)
		}
		if translatedEvents[index].Type == "PR Review Requested" {
			outputBody := outputBody + "\n    URL: " + translatedEvents[index].Payload.PullRequest.URL
			fmt.Println(outputBody)
		}
		if translatedEvents[index].Type == "Repository Created" {
			outputBody := outputBody + "\n    Repository: " + translatedEvents[index].Repo.URL
			fmt.Println(outputBody)
		}
		if translatedEvents[index].Type == "Branch Created" {
			outputBody := outputBody + "\n    Branch: " + translatedEvents[index].Payload.Ref
			fmt.Println(outputBody)
		}
		if translatedEvents[index].Type == "Push" {
			for i := range translatedEvents[index].Payload.Commits {
				outputBody := outputBody + "\n    Commit: " + translatedEvents[index].Payload.Commits[i].URL
				fmt.Println(outputBody)
			}
		}
	}
}

func getEvents (user string, gitToken string) (models.Events, error) {
	eventsEndpoint := "https://api.github.com/users/" + user + "/events"
	req, err := http.NewRequest(http.MethodGet, eventsEndpoint, nil)
	req.Header.Add("Authorization", "Bearer " + gitToken)
	if err != nil {
		return nil, fmt.Errorf("Error creating events request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error executing the request request: %v", err)
	}
	defer res.Body.Close() // Ensure the response body is closed
	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body: %v", err)
	}
	var response models.Events
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON: %v", err)
	}
	return response, nil
}

func translateEvents (events models.Events) (models.Events, error) {
	for index := range events {
		if events[index].Type == "PushEvent" {
			events[index].Type = "Push"
		}
		if events[index].Type == "PullRequestReviewCommentEvent" {
			events[index].Type = "Suggestion Applied"
		}
		if events[index].Type == "PullRequestReviewEvent" {
			events[index].Type = "PR Review Requested"
		}
		if events[index].Type == "CreateEvent" {
			if events[index].Payload.RefType == "branch" {
				events[index].Type = "Branch Created"
			} else if events[index].Payload.RefType == "repository" {
				events[index].Type = "Repository Created"
			}
		}
	}
	return events, nil
}