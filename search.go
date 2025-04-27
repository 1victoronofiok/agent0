package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Search(prompt string, limit int) ([]SearchResponse, error) {
	k, isPresent := os.LookupEnv("SERPER_API_KEY")
	if k == "" || !isPresent {
		return nil, errors.New("")
	}

	url := "https://google.serper.dev/search"
	method := "POST"

	fp := fmt.Sprintf(`{"q":"%s","gl":"ng"}`, prompt)
	payload := strings.NewReader(fp)

	c := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-KEY", k)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var sr []SearchResponse
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}

	return sr, nil

	// TODO: re-map properties from API to fit SearchResponse type
	// return []SearchResponse{
	// 	{
	// 		Title: "The Top AI Conferences To Attend In 2025 - Oxford Abstracts",
	// 		URL:   "https://oxfordabstracts.com/blog/top-ai-conferences-to-attend-in-2024/",
	// 		Order: 1,
	// 	},
	// 	{
	// 		Title: "Top 10 AI Conferences for 2025 | DataCamp",
	// 		URL:   "https://www.datacamp.com/blog/top-ai-conferences",
	// 		Order: 2,
	// 	},
	// 	{
	// 		Title: "AI Conferences 2025: 20 Upcoming Events To Attend",
	// 		URL:   "https://aiconferences.ai/",
	// 		Order: 3,
	// 	},
	// }
}
