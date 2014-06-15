package model

import "time"

type Queue struct {
	Master Master
	Time   time.Time
	Text   string
}

func NewQueueFromPayload(p Payload) (q Queue) {
	return NewQueue(
		p.Master,
		time.Unix(p.Finish, 0),
		p.Text,
	)
}
func NewQueue(masterName string, finish time.Time, text string) Queue {
	return Queue{
		Master{
			Name:       masterName,
			ScreenName: masterName,
		},
		finish,
		text,
	}
}
