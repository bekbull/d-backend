package models

type Item struct {
	Title    string `json:"title"`
	IdNumber string `json:"idNumber"`
	Left     string `json:"left"`
	Image     string `json:"image"`
	Desc     string `json:"desc"`
}