package main

import "time"

// TimerRoutine(time.Second, func() { onSec() }, nil)
func TimerRoutine(tickDur time.Duration, do func(), onErr func(string), run ...bool) chan struct{} {
	if len(run) > 0 && run[0] {
		do()
	}
	quit := make(chan struct{})
	go func() {
		tick := time.NewTicker(tickDur)
		defer func() {
			tick.Stop()
			if onErr != nil {
				if err := recover(); err != nil {
					onErr("")
				}
			}
		}()
		for {
			select {
			case <-tick.C:
				do()
			case <-quit:
				return
			}
		}
	}()
	return quit
}
