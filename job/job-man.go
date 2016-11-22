package job

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/elleFlorio/video-provisioner/job/load"
	"github.com/elleFlorio/video-provisioner/logger"
	"github.com/elleFlorio/video-provisioner/metric"
	"github.com/elleFlorio/video-provisioner/request"
)

var (
	jobs        map[string]request.Request
	mutex_w     = &sync.Mutex{}
	lmbd        float64
	useProfiles = false
	gen         *rand.Rand
)

func init() {
	jobs = make(map[string]request.Request)
	source := rand.NewSource(time.Now().UnixNano())
	gen = rand.New(source)
}

func InitializeJobsManager(lambda float64, probabilities []string) {
	lmbd = lambda

	if len(probabilities) != 0 {
		load.ReadProbabilities(probabilities)
		profiles := load.GetProfilesNames()
		load.ReadProfiles(profiles)
		useProfiles = true
	}

	if useProfiles {
		log.Println("Using load profiles")
	} else {
		log.Printf("Using lambda %f\n", lambda)
	}

}

func ManageJobs(ch_req chan request.Request) {
	log.Println("Started work manager. Waiting for work to do...")
	ch_done := make(chan request.Request)
	var workTime float64
	for {
		select {
		case req := <-ch_req:
			addReqToWorks(req)
			workTime = getWorkTime()
			reqDone := Work(workTime, req, ch_done)
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

func getWorkTime() float64 {
	if !useProfiles {
		return gen.ExpFloat64() * lmbd
	} else {
		return load.GetLoad()
	}
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
