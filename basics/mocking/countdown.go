package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Sleeper interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer, sleeper Sleeper) {
	var err error
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		_, err = fmt.Fprintln(out, i)
		if err != nil {
			fmt.Println(err)
		}
	}

	sleeper.Sleep()
	_, err = fmt.Fprint(out, finalWord)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
	Countdown(os.Stdout, sleeper)
}
