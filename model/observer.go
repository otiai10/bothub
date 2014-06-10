package model

import "fmt"

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
	go func() {
		for {
			q := <-obs.Enq
			fmt.Printf("ここで %+v\n", q)
		}
	}()
}
func (obs *Observer) ListenQueue() {
	go func() {
		for {
			timed := <-obs.Timed
			fmt.Printf("タイムアップ %+v\n", timed)
		}
	}()
}
