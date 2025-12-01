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
	mode = "xterm"
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

	debug("getting a lights string")
	if mode == "ws" {
		if ws, err = ws2811.MakeWS2811(getOpts()); err != nil {
			return
		}

		debug("initiating lights")
		if err = ws.Init(); err != nil {
			return
		}
	}


	lights := makeLights(50, ws)
	showLights(lights, ws)
}

func showLights(lights []*Light, ws *ws2811.WS2811) (err error) {
	debug("starting main loop")
	for { func() {
		mark := time.Now()
		line := " "
		defer func() {
			if mode == "ws" {
				if err = ws.Render(); err != nil {
					return
				}
			} else {
				fmt.Printf("%s\r", line)
				time.Sleep(1 * time.Millisecond)
			}
		}()

		for i, v := range lights {
			var color uint32
			if mark.After(v.SwitchTime) {
				v.Cycle(mark)
			}
			if ws != nil {
				ws.Leds(0)[i] = color
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
		v.StartTime = time.Hour * 17
		v.StopTime = time.Hour * 6
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

	opt.Channels[0].LedCount = 50
	opt.Channels[0].Brightness = 64
	return &opt
}
