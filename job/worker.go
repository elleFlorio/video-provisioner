package job

import (
	"math/rand"
	"time"

	"github.com/elleFlorio/video-provisioner/request"
)

const c_MAXITER = 100

var (
	source rand.Source
	gen    *rand.Rand
)

func init() {
	source = rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

func Work(lambda float64, req request.Request, ch_done chan request.Request) {
	load := gen.ExpFloat64() * lambda
	timer := time.NewTimer(time.Millisecond * time.Duration(load))
	for {
		select {
		case <-timer.C:
			req.ExecTimeMs = computeExecutionTime(req.Start)
			ch_done <- req
			return
		default:
			cpuTest()
		}
	}
}

func cpuTest() float64 {
	plusMinus := false
	pi := 0.0
	for i := 1.0; i < c_MAXITER; i = i + 2.0 {
		if plusMinus {
			pi -= 4.0 / i
			plusMinus = false
		} else {
			pi += 4.0 / i
			plusMinus = true
		}
	}
	return pi
}

func computeExecutionTime(start time.Time) float64 {
	return time.Since(start).Seconds() * 1000
}
