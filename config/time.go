package config

import (
	"github.com/beevik/ntp"
	"time"
)

var EPOCH0 int64 = 0
var EPOCH1 int64 = 7889400000000000

type Time struct {
	Epoch     int
	Current   int64
	StartTime time.Time
	LocalTime time.Time
	LocalTZ   string
}

const ntpServer = "time.apple.com"

var ExactTime *Time

// SetTime will synchronize the ntp server, so we always get a trustable source for UTC time
func SetTime() {
	utime, err := ntp.Time(ntpServer)

	if err != nil {
		panic(err)
	}
	ExactTime = new(Time)
	ExactTime.Current = utime.UTC().UnixNano()
	ExactTime.StartTime = utime.UTC()
	StartedAt = ExactTime.StartTime

	l, _ := time.LoadLocation("Local")
	ExactTime.LocalTZ = l.String()

	go ExactTime.bgTimeSinceStart()

}

func (t *Time) bgTimeSinceStart() {
	ticker := time.NewTicker(time.Nanosecond * 1)

	for range ticker.C {
		SinceStarted = time.Duration(time.Since(StartedAt).Nanoseconds())
		time.Sleep(time.Second * 10)
		if SinceStarted.Nanoseconds() < EPOCH1 {
			t.Epoch = 0
		}
		if SinceStarted.Nanoseconds()%EPOCH1 == 1 {
			t.Epoch = t.Epoch + 1
		}
	}

}

var StartedAt time.Time
var SinceStarted time.Duration

func init() {

	go func() {
		SetTime()
	}()

}
