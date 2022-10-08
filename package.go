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
	server   string
	port     string
	from     string
	password string
	subject  string
	body     string
}
