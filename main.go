package main
import (
	"fmt"
	"os"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func getOpts() (*ws2811.Option) {
	var opt ws2811.Option
	opt.Channels[0].LedCount = 150
	opt.Channels[0].Brightness = 64

	return &opt
}

func main() {
	var err error
	var lights *ws2811.WS2811

	defer func() {
		var ec int
		if lights != nil {
			lights.Fini()
		}
		if err != nil {
			ec = 1
			fmt.Println(err)
		}

		os.Exit(ec)
	}()

	lights, err = ws2811.MakeWS2811(getOpts())
	if err != nil {
		return
	}


}


