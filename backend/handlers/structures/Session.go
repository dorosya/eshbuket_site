package handlers

import "time"

type Session struct {
	Username string
	Expires  time.Time
}

var sessions = make(map[string]Session)
