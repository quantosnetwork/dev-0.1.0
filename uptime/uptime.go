package uptime

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

type Uptime struct {
	State chan bool
	Down  chan bool
}

type UptimeManager interface {
	up() *Uptime
	alive() bool
	update()
	Start()
}

func (u UM) alive() bool {
	return <-u.up().State
}

func (u UM) up() *Uptime {
	return u.UP
}

type UM struct {
	UptimeManager
	UP *Uptime
}

func Manager() UptimeManager {
	um := &UM{}
	up := &Uptime{}
	up.State = make(chan bool)
	up.Down = make(chan bool)
	um.UP = up
	return um
}

func (u UM) update() {
	for alive := true; alive; {
		ticker := time.NewTimer(time.Minute)
		select {
		case <-ticker.C:
			err := ioutil.WriteFile(".uptime", []byte(time.Now().Add(time.Minute).String()), 0644)
			if err != nil {
				panic(err)
			}
		case <-u.UP.Down:
			ticker.Stop()
			alive = false
			err := ioutil.WriteFile(".uptime", []byte{0x00}, 0644)
			if err != nil {
				panic(err)
			}
		}

	}
}

func (u UM) Start() {
	m := Manager()
	go func() {
		for now := range time.Tick(time.Second) {
			m.update()
			fmt.Println(now)
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			close(m.up().State)
			<-m.up().Down
			os.Exit(1)
		}
	}()
}
