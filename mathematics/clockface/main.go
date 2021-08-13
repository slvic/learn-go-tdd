package main

import (
	"os"
	"time"

	clockface "github.com/vkenrik117/learn-go-tdd/mathematics"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
