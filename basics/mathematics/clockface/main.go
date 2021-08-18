package main

import (
	clockface "github.com/vkenrik117/learn-go-tdd/basics/mathematics"
	"os"
	"time"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
