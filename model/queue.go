package model

import "time"

type Queue struct {
	Master Master
	Time   time.Time
	Text   string
}

type Observer struct {
	enq *<-chan Queue
	deq *chan<- Queue
}
