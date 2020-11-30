package main

import (
	"flag"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"

	"github.com/jrzhag12/server-monitor/pkg/loader"
	"github.com/jrzhag12/server-monitor/pkg/terminal"
)

func main() {
	url := flag.String("url", "", "/debug/vars")
	flag.Parse()
	if len(*url) == 0 {
		flag.Usage()
		return
	}

	l := loader.NewMemStatsLoader(*url)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	controller := terminal.NewController()

	event := ui.PollEvents()
	tick := time.Tick(time.Second)

	for {
		select {
		case e := <-event:
			switch e.Type {
			case ui.KeyboardEvent:
				// quit on any keyboard event
				return
			case ui.ResizeEvent:
				controller.Resize()
			}
		case <-tick:
			stat, err := l.Load()
			if err != nil {
				log.Println(err)
				break
			}
			// update dashboard every second
			controller.Render(stat)
		}
	}
}
