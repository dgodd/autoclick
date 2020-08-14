package main

import (
	"fmt"
	"time"
	"sync/atomic"
	"github.com/robotn/gohook"
	"github.com/go-vgo/robotgo"
	"github.com/martinlindhe/notify"
)

func main() {
	// fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	// hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
	// 	fmt.Println("ctrl-shift-q")
	// 	hook.End()
	// })

	var count int64
	var click bool
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if click {
				robotgo.Click()
			}
		}
	}()

	fmt.Println("--- Please press c 4 times to toggle---")
	notify.Notify("autoclick", "notice", "press c 4 times to toggle", "")
	hook.Register(hook.KeyDown, []string{"c"}, func(e hook.Event) {
		newCount := atomic.AddInt64(&count, 1)
		go func() {
			time.Sleep(1*time.Second)
			newCount := atomic.AddInt64(&count, -1)
			if newCount < 0 {
				atomic.StoreInt64(&count, 0)
			}
		}()
		if newCount >= 4 {
			atomic.StoreInt64(&count, 0)
			click = !click
			fmt.Println("CLICK:", click)
			if click {
				notify.Notify("autoclick", "notice", "Turn autoclick ON", "")
			}	else{
				notify.Notify("autoclick", "notice", "Turn autoclick OFF", "")
			}
		}
		fmt.Println("pressed: c", newCount)
	})

	s := hook.Start()
	<-hook.Process(s)
}
