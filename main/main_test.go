package fake_main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

const calls = 100

func TestMainRun(t *testing.T) {
	t.Parallel()
	times := make([]time.Duration, calls)
	for i := 0; i < calls; i++ {
		cmd := exec.Command("go", "run", "main.go")

		now := time.Now()
		if err := cmd.Start(); err != nil {
			if err != nil {
				t.Error(err)
			}
		}

		err := cmd.Wait()
		if err != nil {
			t.Error(err)
		}

		times = append(times, time.Since(now))
	}

	var sum int64
	for _, duration := range times {
		sum += duration.Milliseconds()
	}

	fmt.Println("avg:", sum/calls)
}
