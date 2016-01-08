package main

import "time"

type Event struct {
	Client      string
	User        string
	Identifier  string
	Time        time.Time
	Method      string
	Path        string
	Protocol    string
	StatusCode  int
	PayloadSize int
	UserAgent   string
	Referer     string
	Host        string
}
