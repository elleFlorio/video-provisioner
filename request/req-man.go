package request

import (
	"log"
	"net/http"
	_ "runtime"
	_ "sync"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
	"github.com/elleFlorio/video-provisioner/logger"
	"github.com/elleFlorio/video-provisioner/network"
	"github.com/elleFlorio/video-provisioner/utils"
)

var (
	//requests map[string]Request
	//mutex_r  = &sync.Mutex{}
	ch_req  chan struct{}
	counter int
)

func init() {
	//requests = make(map[string]Request)
	ch_req = make(chan struct{})
	counter = 0
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

	updateReqCounter()

	return req, nil
}

// func GetRequest(reqId string) (Request, bool) {
// 	req, ok := requests[reqId]
// 	return req, ok
// }

func IsServiceWaiting() bool {
	// mutex_r.Lock()
	// requestsPending := len(requests)
	// mutex_r.Unlock()

	// if requestsPending != 0 {
	// 	return true
	// }

	// return false

	return false
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
			//addRequestToHistory(reqDone)
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
				//addRequestToHistory(reqDone)
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

func updateReqCounter() {
	ch_req <- struct{}{}
}

func startReqCounter() {
	ticker := time.NewTicker(time.Duration(60) * time.Second)

	for {
		select {
		case <-ticker.C:
			logger.LogRequestsPerMinute(counter)
			counter = 0
		case <-ch_req:
			counter += 1
		}
	}
}
