package models

type Post struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Domain   string `json:"domain"`
	Author   string `json:"author"`
	PostedOn string `json:"postedOn"`
}