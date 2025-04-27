package main

type SearchResponse struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Order int    `json:"order,omitempty"`
}

type PageDetail struct {
	SearchResponse
	Content string `json:"content"`
}
