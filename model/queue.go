package model

import "time"

type Queue struct {
	Master Master
	Time   time.Time
	Text   string
}
