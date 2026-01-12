package handlers

import (
	"time"
)

type Session struct {
	Username string
	Expires  time.Time
}

var Sessions = make(map[string]Session)
