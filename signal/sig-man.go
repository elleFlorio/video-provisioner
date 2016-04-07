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
	if useDiscovery {
		ch_stop <- struct{}{}
		log.Println("Stopped keep alive goroutine")
		discovery.UnregisterFromEtcd()
		log.Println("Unregistered from etcd")
	}
	for job.IsServiceWorking() {
		log.Println("Waiting for jobs to complete...")
		time.Sleep(time.Duration(1) * time.Second)
	}
	for request.IsServiceWaiting() {
		log.Println("Waiting for responses to requests...")
		time.Sleep(time.Duration(1) * time.Second)
	}
	log.Fatalln("Done. Shutting down")
}
