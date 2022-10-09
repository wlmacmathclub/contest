package main

type User struct {
	email         string
	firstName     string
	lastName      string
	firstTeacher  string
	secondTeacher string
}

type Contest struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type MailConfig struct {
	Email      string `json:"email"`
	Name       string `json:"name"`
	Publickey  string `json:"publickey"`
	Privatekey string `json:"privatekey"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
}

type JsonData struct {
	Contest Contest    `json:"contest"`
	Config  MailConfig `json:"config"`
}
