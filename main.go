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
	gitToken := os.Getenv("GITHUB_TOKEN")
	fmt.Println(getEvents(user, gitToken))
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