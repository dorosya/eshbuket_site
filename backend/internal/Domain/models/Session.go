package models

import (
	"time"
)

type Session struct {
	Username string
	Expires  time.Time
}
