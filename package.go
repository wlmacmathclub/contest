package main

import "time"

type User struct {
	email         string
	firstName     string
	lastName      string
	firstTeacher  string
	secondTeacher string
}

type Contest struct {
	name string
	date time.Time
}
