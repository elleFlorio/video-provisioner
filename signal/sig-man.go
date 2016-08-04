package signal

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
	"github.com/elleFlorio/video-provisioner/job"
	"github.com/elleFlorio/video-provisioner/request"
)

const c_WAIT_COUNTER_LIMIT = 30

func MonitorSignals(ch_stop chan struct{}, useDiscovery bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			go shutDown(ch_stop, useDiscovery)
		}
	}
}

func shutDown(ch_stop chan struct{}, useDiscovery bool) {
	log.Println("Received shutdown signal")
	waitCounter := 0
	// Unregister if needed
	if useDiscovery {
		ch_stop <- struct{}{}
		log.Println("Stopped keep alive goroutine")
		discovery.UnregisterFromEtcd()
		log.Println("Unregistered from etcd")
	}
	// Complete current jobs
	for job.IsServiceWorking() {
		log.Println("Waiting for jobs to complete...")
		time.Sleep(time.Duration(1) * time.Second)
	}
	// Wait for responses if its needed
	for request.IsServiceWaiting() && waitCounter < c_WAIT_COUNTER_LIMIT {
		log.Println("Waiting for responses to requests...")
		waitCounter += 1
		time.Sleep(time.Duration(1) * time.Second)
	}
	// Stop req counter and log the requests arrived till now
	request.StopReqCounter()
	time.Sleep(time.Duration(1) * time.Second)
	log.Println("Logged requests arrived")

	log.Fatalln("Done. Shutting down")
}
