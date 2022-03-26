package config

import (
	"github.com/beevik/ntp"
	"sync"
	"time"
)

const (
	NANOSECONDS_PER_EPOCH int64 = 7889400000000000 // 1 EPOCH = 3 months there is 4 EPOCHS per year

)

var EPOCH0 int64 = 0
var EPOCH1 int64 = 7889400000000000

type Time struct {
	Epoch                int
	Current              int64
	StartTime            time.Time
	LocalTime            time.Time
	LocalTZ              string
	StartTimeInitialized bool
}

const ntpServer = "time.apple.com"

var ExactTime *Time

// SetTime will synchronize the ntp server, so we always get a trustable source for UTC time
func SetTime() {
	utime, err := ntp.Time(ntpServer)
	if err != nil {
		panic(err)
	}

	ExactTime.Current = utime.UTC().UnixNano()
	if ExactTime.StartTimeInitialized == false {
		ExactTime.StartTime = utime.UTC()
		ExactTime.StartTimeInitialized = true
	} else {
		timeSinceStart := time.Now().Sub(ExactTime.StartTime)
		tns := timeSinceStart.Nanoseconds()
		if tns < EPOCH1 {
			ExactTime.Epoch = 0
		} else {
			epoch := 0
			if tns%EPOCH1 == 1 {
				epoch = ExactTime.Epoch + 1
				ExactTime.Epoch = epoch
			}
		}
	}

	l, _ := time.LoadLocation("Local")
	ExactTime.LocalTZ = l.String()

}

// SetStartTime + update UTC with ntp server at intervals
func SetStartTime() {
	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(2)
	go SetTime()

	go func() {
		for {
			SetTime()
		}
	}()
	wg.Wait()

}
