package main
import (
	"fmt"
	"os"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)
func mayn() {
	var err error
	var lights *ws2811.WS2811

	debug("starting program")

	defer func() {
		debug("dying")
		var ec int
		if lights != nil {
			lights.Fini()
		} else {
			debug("obj was nil")
		}
		if err != nil {
			ec = 1
			fmt.Println(err)
		}

		os.Exit(ec)
	}()

	debug("getting a lights string")
	lights, err = ws2811.MakeWS2811(getOpts())
	if err != nil {
		return
	}

	debug("initiating lights")
	err = lights.Init()
	if err != nil {
		return
	}
	debug("examining light object")
	x := lights.Leds(0)
	debug("%+#v", x)

}
