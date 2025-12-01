package main
import (
	"time"
)

type Light struct {
	On                      bool           // Defines whether or not the light is switched on at the time of updating
	Color                   int           // Index of a slice listing colors that will be displayed
	WsColor                 uint32
	SwitchTime              time.Time     // The Clock time at which the light should flip the On boolean value
	OnDuration, OffDuration time.Duration // The duration after a flip that the light should remain in the new state
	StartTime, EndTime      time.Duration // The duration after 00:00:00 local time that the light should start its cycle or end its cycle
}

func (l *Light) SetDurations() {
	l.OnDuration = duration()
	l.OffDuration = duration()
}

func (l *Light) Cycle(mark, midnight time.Time) () {
	var t time.Duration
	if mark.After(midnight.Add(l.StartTime)) || mark.Before(midnight.Add(l.EndTime)) {
		if l.On {
			l.WsColor = WsOff
			t = l.OffDuration
			l.On = false
		} else {
			t = l.OnDuration
			l.On = true
			l.WsColor = WsColorWheel[l.Color]
		}
		l.SwitchTime = mark.Add(t)
	} else {
		l.On = false
		l.WsColor = WsOff
		l.SwitchTime = midnight.Add(l.StartTime)
	}
}

func (l *Light)GetColor() string {
	if l.On {
		switch mode {
		case "xterm":
			return l.XtermGetColor()
		case "tty":
			return l.TTYGetColor()
		default:
		}
	}

	return " "
}

func (l *Light) XtermGetColor() string {
	return XtermColorWheel[l.Color] + "█" + XtermReset
}

func (l *Light) TTYGetColor() string {
	return TTYColorWheel[l.Color]
}
