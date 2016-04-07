package job

import (
	"log"
	"runtime"
	"sync"

	"github.com/elleFlorio/video-provisioner/logger"
	"github.com/elleFlorio/video-provisioner/metric"
	"github.com/elleFlorio/video-provisioner/request"
)

var (
	jobs    map[string]request.Request
	mutex_w = &sync.Mutex{}
)

func init() {
	jobs = make(map[string]request.Request)
}

func ManageJobs(ch_req chan request.Request) {
	log.Println("Started work manager. Waiting for work to do...")
	ch_done := make(chan request.Request)
	for {
		select {
		case req := <-ch_req:
			addReqToWorks(req)
			go Work(getLambda(), req, ch_done)
		case reqDone := <-ch_done:
			logger.LogExecutionTime(reqDone.ExecTimeMs)
			request.FinalizeReq(reqDone)
			removeReqFromWorks(reqDone.ID)
			metric.SendExecutionTime(reqDone.ExecTimeMs)
		}
	}
}

func addReqToWorks(req request.Request) {
	mutex_w.Lock()
	jobs[req.ID] = req
	mutex_w.Unlock()
	runtime.Gosched()
}

func removeReqFromWorks(id string) {
	mutex_w.Lock()
	delete(jobs, id)
	mutex_w.Unlock()
	runtime.Gosched()
}

func getLambda() float64 {
	return 0.0
}

func IsServiceWorking() bool {
	defer runtime.Gosched()
	mutex_w.Lock()
	jobsInProgress := len(jobs)
	mutex_w.Unlock()

	if jobsInProgress != 0 {
		return true
	}

	return false
}
