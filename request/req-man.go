package request

import (
	"log"
	"net/http"
	_ "runtime"
	_ "sync"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
	"github.com/elleFlorio/video-provisioner/network"
	"github.com/elleFlorio/video-provisioner/utils"
)

var (
	//requests map[string]Request
	//mutex_r  = &sync.Mutex{}
	counter = 1
)

func init() {
	//requests = make(map[string]Request)
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

// func addRequestToHistory(req Request) {
// 	mutex_r.Lock()
// 	requests[req.ID] = req
// 	mutex_r.Unlock()
// 	runtime.Gosched()
// }

// func UpdateRequestInHistory(reqId string) bool {
// 	deleted := false
// 	mutex_r.Lock()
// 	req := requests[reqId]
// 	req.Counter -= 1
// 	if req.Counter <= 0 {
// 		delete(requests, reqId)
// 		deleted = true
// 	} else {
// 		requests[reqId] = req
// 	}
// 	mutex_r.Unlock()
// 	runtime.Gosched()
// 	return deleted
// }
