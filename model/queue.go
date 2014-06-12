package model

import "time"

type Queue struct {
	Master Master
	Time   time.Time
	Text   string
}

func NewQueueFromPayload(p Payload) (q Queue) {
	_dur, _ := time.ParseDuration("3s")
	return NewQueue("otiai10", time.Now().Add(_dur), "3秒後です"+time.Now().Format("2006010215040506"))
}
func NewQueue(masterName string, finish time.Time, text string) Queue {
	return Queue{
		Master{
			Name: masterName,
		},
		finish,
		text,
	}
}
