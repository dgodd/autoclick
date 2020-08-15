package main

import (
	"fmt"
	"time"
	"sync/atomic"
	"os"
	"log"
    "path/filepath"
	"github.com/robotn/gohook"
	"github.com/go-vgo/robotgo"
	"github.com/martinlindhe/notify"
	"github.com/getlantern/systray"
)

func main() {
	onExit := func() {
		fmt.Println("BYE")
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DIR:", dir)
	iconFile := filepath.Join(dir, "..", "..", "Resources", "click-icon.png")

	systray.SetTitle("Auto Click")
	systray.SetTooltip("Click 'c' 4 times")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	listenForC(iconFile)
}

func listenForC(iconFile string) {
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
	notify.Notify("autoclick", "notice", "press c 4 times to toggle", iconFile)
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
				notify.Notify("autoclick", "notice", "Turn autoclick ON", iconFile)
			}	else{
				notify.Notify("autoclick", "notice", "Turn autoclick OFF", iconFile)
			}
		}
		fmt.Println("pressed: c", newCount)
	})

	s := hook.Start()
	<-hook.Process(s)
}
