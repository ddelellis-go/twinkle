package main

import (
	"fmt"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"flag"
	"math/rand"
	"os"
	"time"
	"github.com/ddelellis-pkg/debugger"
)

var mode string
var count int

func main() {

	flag.StringVar(&mode, "displaymode", "ws", "'xterm', 'tty', or 'ws'")
	flag.IntVar(&count, "bulbcount", 20, "number of bulbs to operate on")
	flag.Parse()

	debugger.Verbose = true
	var err error
	var ws *ws2811.WS2811
	var lights []*Light

	debug("starting program")

	defer func() {
		os.Exit(shutdown(ws, err))
	}()

	ws, lights, err = initLights(count)

	showLights(lights, ws)
}

func initLights(ledCount int) (ws *ws2811.WS2811, lights []*Light, err error) {
	if mode == "ws" {
		debug("getting a WS2811 obj")
		if ws, err = ws2811.MakeWS2811(getOpts(count)); err != nil {
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

	lights = makeLights(count)

	return
}

func showLights(lights []*Light, ws *ws2811.WS2811) (err error) {
	debug("starting main loop")


	for { func() {
		mark := time.Now()
		var nextCycle time.Time
		midnight := time.Date(mark.Year(), mark.Month(), mark.Day(), 0, 0, 0, 0,time.Local)
		line := " "
		defer func() {
			if mode == "ws" {
				if err = ws.Render(); err != nil {
					return
				}
				if err = ws.Wait(); err != nil {
					return
				}
			}
			sleepDuration := time.Until(nextCycle)
			fmt.Printf("%s\r",line)
			time.Sleep(sleepDuration)
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
			if nextCycle.IsZero() {
				nextCycle = v.SwitchTime
			} else {
				if v.SwitchTime.Before(nextCycle) {
					nextCycle = v.SwitchTime
				}
			}
			lights[i] = v
		}
	}()}
}

func makeLights(count int) (lights []*Light) {
	lights = make([]*Light, count)
	mark := time.Now()
	for i, _ := range lights {
		var v Light
		v.SwitchTime = mark
		v.SetDurations()
		v.Color = i % 4 + 1
		v.WsColor = WsColorWheel[ v.Color ]
		v.StartTime = time.Hour * 17
		v.EndTime = time.Hour * 6
		lights[i] = &v
	}

	return lights
}

func debug(layout string, args ...any) {
	debugger.Debug(layout, args...)
}

func duration() (t time.Duration) {
	x := 750 + int64(2000.0*rand.ExpFloat64())
	if x >= 5000 {
		x = 1500 + rand.Int63n(3500)
	}
	t = time.Millisecond * time.Duration(x)
	return
}

func getOpts(ledCount int) *ws2811.Option {
	var opt ws2811.Option
	opt = ws2811.DefaultOptions

	opt.Channels[0].LedCount = ledCount
	opt.Channels[0].Brightness = 255
	return &opt
}

func shutdown(ws *ws2811.WS2811, err error) (exitCode int) {
	exitCode = debugger.DumpErrorStack(err)
	if ws != nil {
		ws.Fini()
	}
	return
}
