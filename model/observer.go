package model

import "fmt"
import "time"

var instance_observer *Observer

type Observer struct {
	Enq   chan Queue
	Timed chan Queue
}

func GetObserver() *Observer {
	if instance_observer == nil {
		instance_observer = &Observer{
			make(chan Queue),
			make(chan Queue),
		}
	}
	return instance_observer
}
func (obs *Observer) Enqueue(q Queue) (e error) {
	obs.Enq <- q
	return
}
func (obs *Observer) WaitQueue() {
	for {
		q := <-obs.Enq
		obs.generateRoutine(q)
	}
}
func (obs *Observer) ListenQueue() {
	for {
		timed := <-obs.Timed
		fmt.Printf("タイムアップ %+v\n", timed)
	}
}
func (obs *Observer) generateRoutine(q Queue) {
	go func() {
		dur := q.Time.Sub(time.Now())
		<-time.After(dur)
		obs.Timed <- q
	}()
}
