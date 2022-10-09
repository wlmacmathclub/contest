package main

type User struct {
	email         string
	firstName     string
	lastName      string
	firstTeacher  string
	secondTeacher string
}

type Contest struct {
	name string `json:"name"`
	date string `json:"date"`
}

type MailConfig struct {
	email      string `json:"email"`
	name       string `json:"name"`
	publickey  string `json:"publickey"`
	privatekey string `json:"privatekey"`
	subject    string `json:"subject"`
	body       string `json:"body"`
}

type JsonData struct {
	contest Contest    `json:"contest"`
	config  MailConfig `json:"config"`
}
