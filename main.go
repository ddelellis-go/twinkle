package main

import (
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	mode = "ws"
)

func main() {
	var err error
	var ws *ws2811.WS2811

	debug("starting program")

	defer func() {
		debug("dying")
		var ec int
		if ws != nil {
			ws.Fini()
		} else {
			err = fmt.Errorf("ws object was nil")
		}
		if err != nil {
			ec = 1
			fmt.Println(err)
		}

		os.Exit(ec)
	}()

	if mode == "ws" {
		debug("getting a WS2811 obj")
		if ws, err = ws2811.MakeWS2811(getOpts()); err != nil {
			return
		}
		if ws == nil {
			err = fmt.Errorf("returned poiner was nil")
			return
		}

		debug("initiating lights")
		if err = ws.Init(); err != nil {
			return
		}
	}


	debug("getting a lights string")
	lights := makeLights(150, ws)
	showLights(lights, ws)
}

func showLights(lights []*Light, ws *ws2811.WS2811) (err error) {
	debug("starting main loop")
	for { func() {
		mark := time.Now()
		midnight := time.Date(mark.Year(), mark.Month(), mark.Day(), 0, 0, 0, 0,time.Local)
		line := " "
		defer func() {
			if mode == "ws" {
				if err = ws.Render(); err != nil {
					return
				}
				fmt.Printf("%s\r", line)
				if err = ws.Wait(); err != nil {
					return
				}
			}
		}()

		for i, v := range lights {
			if mark.After(v.SwitchTime) {
				v.Cycle(mark, midnight)
			}
			if ws != nil {
				ws.Leds(0)[i] = v.WsColor
			}
			if mode == "xterm" || mode == "tty" {
				line += v.GetColor()
			}
			lights[i] = v
		}
	}()}
}

func makeLights(count int, ws *ws2811.WS2811) (lights []*Light) {
	lights = make([]*Light, count)
	mark := time.Now()
	for i, _ := range lights {
		var v Light
		v.SwitchTime = mark
		v.SetDurations()
		v.Color = i%4 + 1
		v.WsColor = WsColorWheel[ v.Color ]
		v.StartTime = time.Hour * 17
		v.EndTime = time.Hour * 6
		lights[i] = &v
	}

	return lights
}

func debug(layout string, args ...any) {
	if !strings.HasSuffix(layout, "\n") {
		layout += "\n"
	}
	fmt.Printf(layout, args...)
}

func duration() (t time.Duration) {
	x := 750 + int64(2000.0*rand.ExpFloat64())
	if x >= 5000 {
		x = 1500 + rand.Int63n(3500)
	}
	t = time.Millisecond * time.Duration(x)
	return
}

func getOpts() *ws2811.Option {
	var opt ws2811.Option
	opt = ws2811.DefaultOptions

	opt.Channels[0].LedCount = 150
	opt.Channels[0].Brightness = 255
	debug("%+#v", opt)
	return &opt
}
