package main
import (
	"math/rand"
	"fmt"
	"strings"
	"os"
	"time"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	Reset = "\033[0m"
	Red = "\033[31m"
	Green = "\033[32m"
	Yellow = "\033[33m"
	Blue = "\033[34m"
	Magenta = "\033[35m"
	Cyan = "\033[36m"
	Gray = "\033[37m"
	White = "\033[97m"
)

var ColorWheel = []string{Reset, Red, Green, Yellow, Blue}

func(l *Light)ShowColor() string {
	return ColorWheel[l.Color]+"0"+Reset
}

func main() {
	var err error

	debug("starting program")

	defer func() {
		debug("dying")
		var ec int
		if err != nil {
			ec = 1
			fmt.Println(err)
		}

		os.Exit(ec)
	}()

	lights := makeLights(80)

	for {
		line := " "
		mark := time.Now()

		for i,v := range lights {
			if mark.After(v.SwitchTime) {
				v.Cycle(mark)
			}
			if v.On {
				line += v.ShowColor()
			} else {
				line += fmt.Sprintf(" ")
			}
			lights[i] = v
		}
		time.Sleep(17 * time.Millisecond)
		fmt.Printf("%s\r",line)
	}

}

func makeLights(count int) (lights []Light) {
	lights = make([]Light, count)
	mark := time.Now()
	for i, v := range lights {
		v.SwitchTime = mark
		v.SetDurations()
		v.Color = i % 4 + 1
		lights[i] = v
	}

	return lights
}

func (l *Light)Cycle(mark time.Time) {
	var t time.Duration
	if l.On {
		t = l.OffDuration
		l.On = false
	} else {
		t = l.OnDuration
		l.On = true
	}
	l.SwitchTime = mark.Add(t)
}


func (l *Light) SetDurations() {
	l.OnDuration  = duration()
	l.OffDuration = duration()
}

type Light struct {
	On bool
	Color int
	SwitchTime time.Time
	OnDuration, OffDuration time.Duration
}

func debug(layout string, args ...any) {
	if !strings.HasSuffix(layout,"\n") {
		layout += "\n"
	}
	fmt.Printf(layout, args...)
}

func duration() (t time.Duration) {
	x := 750+int64(2000.0*rand.ExpFloat64())
	if x >= 5000 {
		x = 1500 + rand.Int63n(3500)
	}
	t = time.Millisecond * time.Duration(x)
	return
}

func getOpts() (*ws2811.Option) {
	var opt ws2811.Option
	opt = ws2811.DefaultOptions

	opt.Channels[0].LedCount = 80
	opt.Channels[0].Brightness = 64
	return &opt
}
