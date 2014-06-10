package model

import "time"

type Queue struct {
	Master Master
	Time   time.Time
	Text   string
}

func NewQueueFromPayload(p Payload) (q Queue) {
	_dur, _ := time.ParseDuration("10s")
	return NewQueue("otiai10です", time.Now().Add(_dur), "ほげえええ")
}
func NewQueue(masterName string, finish time.Time, text string) Queue {
	return Queue{
		Master{
			masterName,
		},
		finish,
		text,
	}
}
