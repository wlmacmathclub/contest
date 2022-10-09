package main

type User struct {
	email         string
	firstName     string
	lastName      string
	firstTeacher  string
	secondTeacher string
}

type Contest struct {
	name string
	date string
}

type MailConfig struct {
	email      string
	name       string
	publickey  string
	privatekey string
	subject    string
	body       string
}
