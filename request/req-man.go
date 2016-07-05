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
	//requests map[string]Request
	//mutex_r  = &sync.Mutex{}
	reqArr       []string
	ch_req_arr   chan string
	ch_req_done  chan struct{}
	counter_done int
)

func init() {
	//requests = make(map[string]Request)
	reqArr = []string{}
	ch_req_arr = make(chan string)
	ch_req_done = make(chan struct{})
	counter_done = 0
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

	updateReqArr(requestID)

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

	if contains(reqArr, reqDone.ID) {
		updateReqDoneCounter()
	}
}

func StartReqCounter() {
	go startReqCounter()
}

func updateReqArr(id string) {
	ch_req_arr <- id
}

func updateReqDoneCounter() {
	ch_req_done <- struct{}{}
}

func startReqCounter() {
	ticker := time.NewTicker(time.Duration(60) * time.Second)

	for {
		select {
		case <-ticker.C:
			logger.LogRequestsArrivedPerMinute(len(reqArr))
			logger.LogRequestsDonePerMinute(counter_done)
			metric.SendRequestsArrived(len(reqArr))
			metric.SendRequestsDone(counter_done)
			reqArr = reqArr[:0]
			counter_done = 0
		case id := <-ch_req_arr:
			reqArr = append(reqArr, id)
		case <-ch_req_done:
			counter_done += 1
		}
	}
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}

	return false
}
