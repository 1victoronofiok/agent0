package main

func Search(query string, limit int) []SearchResponse {

	// TODO: re-map properties from API to fit SearchResponse type
	return []SearchResponse{
		{
			Title: "The Top AI Conferences To Attend In 2025 - Oxford Abstracts",
			URL:   "https://oxfordabstracts.com/blog/top-ai-conferences-to-attend-in-2024/",
			Order: 1,
		},
		{
			Title: "Top 10 AI Conferences for 2025 | DataCamp",
			URL:   "https://www.datacamp.com/blog/top-ai-conferences",
			Order: 2,
		},
		{
			Title: "AI Conferences 2025: 20 Upcoming Events To Attend",
			URL:   "https://aiconferences.ai/",
			Order: 3,
		},
	}
}
