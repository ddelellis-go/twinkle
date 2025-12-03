package main
import (
)

const (
	TTYOff     = " "
	TTYRed     = "R"
	TTYGreen   = "G"
	TTYBlue    = "B"
	TTYYellow  = "Y"
	TTYMagenta = "M"
	TTYCyan    = "C"
	TTYGray    = "A"
	TTYWhite   = "W"

	XtermReset   = "\033[0m"
	XtermRed     = "\033[31m"
	XtermGreen   = "\033[32m"
	XtermBlue    = "\033[94m"
	XtermYellow  = "\033[33m"
	XtermMagenta = "\033[35m"
	XtermCyan    = "\033[36m"
	XtermGray    = "\033[37m"
	XtermWhite   = "\033[97m"

	WsOff     = uint32(0x000000)
	WsRed     = uint32(0xff0000)
	WsGreen   = uint32(0x00ff00)
	WsBlue    = uint32(0x1e90ff)
	WsYellow  = uint32(0xffff00)
	WsMagenta = uint32(0xff00ff)
	WsCyan    = uint32(0x00ffff)
	WsGray    = uint32(0x7f7f7f)
	WsWhite   = uint32(0xffffff)
)

var (
	XtermColorWheel = []string{XtermReset, XtermRed, XtermGreen, XtermYellow, XtermBlue}
	TTYColorWheel = []string{TTYOff, TTYRed, TTYGreen, TTYYellow, TTYBlue}
	WsColorWheel = []uint32{WsOff, WsRed, WsGreen, WsYellow, WsBlue}
)
