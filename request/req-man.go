package request

import (
	"log"
	"net/http"
	_ "runtime"
	_ "sync"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
	"github.com/elleFlorio/video-provisioner/logger"
	"github.com/elleFlorio/video-provisioner/metric"
	"github.com/elleFlorio/video-provisioner/network"
	"github.com/elleFlorio/video-provisioner/utils"
)

var (
	ch_req_arr  chan struct{}
	ch_req_stop chan struct{}
	counter_arr int
)

func init() {
	ch_req_arr = make(chan struct{})
	ch_req_stop = make(chan struct{})
	counter_arr = 0
}

func CreateReq(r *http.Request) (Request, error) {
	var err error
	var requestID string
	var start = time.Now()

	//read request
	message, err := network.ReadMessage(r)
	if err != nil {
		log.Println("Cannot read message")
		return Request{}, err
	}
	requestID = message.Args
	if requestID == "" {
		log.Println("New request, generating ID")
		requestID, _ = utils.GenerateUUID()
	}
	log.Printf("Received request %s from %s\n", requestID, message.Sender)

	toService, _ := network.ReadParam("service", r)

	// Counter is not used, but it can be useful if I will need to
	// implement requests with multiple destinations
	req := Request{
		ID:         requestID,
		From:       message.Sender,
		To:         toService,
		Counter:    1,
		Start:      start,
		ExecTimeMs: 0,
	}

	updateReqArr()

	return req, nil
}

func FinalizeReq(reqDone Request) {
	var err error
	isEndpointSet := network.IsEndpointSet()
	if reqDone.To != "" {
		err = network.SendMessageToSpecificService(reqDone.ID, reqDone.To)
		if err != nil {
			log.Println("Cannot dispatch message to service", reqDone.To)
			return
		}
		if !isEndpointSet {
			discovery.AddRequestToHistory(reqDone.ID, reqDone.Start)

		}
	} else {
		if network.GetDestinationsNumber() > 0 {
			err = network.SendMessageToDestination(reqDone.ID)
			if err != nil {
				log.Println(err)
				if isEndpointSet {
					network.RespondeToEndpoint(reqDone.ID, "done")
				} else {
					network.RespondeToRequest(reqDone.From, reqDone.ID, "done")
				}
				return
			}
			if !isEndpointSet {
				discovery.AddRequestToHistory(reqDone.ID, reqDone.Start)
			}
		} else {
			if isEndpointSet {
				network.RespondeToEndpoint(reqDone.ID, "done")
			} else {
				network.RespondeToRequest(reqDone.From, reqDone.ID, "done")
			}
		}
	}
}

func StartReqCounter() {
	go startReqCounter()
}

func StopReqCounter() {
	ch_req_stop <- struct{}{}
}
func updateReqArr() {
	ch_req_arr <- struct{}{}
}

func startReqCounter() {
	ticker := time.NewTicker(time.Duration(60) * time.Second)

	for {
		select {
		case <-ticker.C:
			logger.LogRequestsArrivedPerMinute(counter_arr)
			metric.SendRequestsArrived(counter_arr)
			counter_arr = 0
		case <-ch_req_arr:
			counter_arr += 1
		case <-ch_req_stop:
			logger.LogRequestsArrivedPerMinute(counter_arr)
			metric.SendRequestsArrived(counter_arr)
			counter_arr = 0
			return
		}
	}
}
