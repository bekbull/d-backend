package models

type Item struct {
	Title    string `json:"title"`
	IdNumber int16  `json:"idNumber"`
	Left     int8   `json:"left"`
	Image    string `json:"image"`
	Desc     string `json:"desc"`
}
