package request

import (
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/elleFlorio/video-provisioner/network"
)

var (
	requests map[string]Request
	mutex_c  = &sync.Mutex{}
	mutex_r  = &sync.Mutex{}
	counter  = 1
)

func init() {
	requests = make(map[string]Request)
}

func CreateReq(r *http.Request, nDestination int) (Request, error) {
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
		requestID = strconv.Itoa(readAndIncrementCounter())
	}
	log.Printf("Received request %s from %s\n", requestID, message.Sender)

	toService, _ := network.ReadParam("service", r)

	req := Request{
		ID:         requestID,
		From:       message.Sender,
		To:         toService,
		Counter:    nDestination,
		Start:      start,
		ExecTimeMs: 0,
	}

	return req, nil
}

func readAndIncrementCounter() int {
	mutex_c.Lock()
	c := counter
	counter++
	mutex_c.Unlock()
	runtime.Gosched()

	return c
}

func GetRequest(reqId string) (Request, bool) {
	req, ok := requests[reqId]
	return req, ok
}

func IsServiceWaiting() bool {
	mutex_r.Lock()
	requestsPending := len(requests)
	mutex_r.Unlock()

	if requestsPending != 0 {
		return true
	}

	return false
}

func FinalizeReq(reqDone Request, destinations []string) {
	if reqDone.To != "" {
		err := network.SendMessageToSpecificService(reqDone.ID, reqDone.To)
		if err != nil {
			log.Println("Cannot dispatch message to service", reqDone.To)
			return
		}
		addRequestToHistory(reqDone)
	} else {
		if len(destinations) > 0 {
			errCounter := network.SendMessageToDestinations(reqDone.ID, destinations)
			if errCounter < len(destinations) {
				// This is for requests to multiple destinations
				// because I have to wait till every destination
				// responde me before consider the request complete
				reqCounter := len(destinations) - errCounter
				reqDone.Counter = reqCounter
				addRequestToHistory(reqDone)
			}
			if errCounter > 0 {
				log.Println("Cannot dispatch message to all the destinations")
				if errCounter == len(destinations) {
					network.RespondeToRequest(reqDone.From, reqDone.ID, "done")
				}
				return
			}
		} else {
			network.RespondeToRequest(reqDone.From, reqDone.ID, "done")
		}
	}
}

func addRequestToHistory(req Request) {
	mutex_r.Lock()
	requests[req.ID] = req
	mutex_r.Unlock()
	runtime.Gosched()
}

func UpdateRequestInHistory(reqId string) bool {
	deleted := false
	mutex_r.Lock()
	req := requests[reqId]
	req.Counter -= 1
	if req.Counter <= 0 {
		delete(requests, reqId)
		deleted = true
	} else {
		requests[reqId] = req
	}
	mutex_r.Unlock()
	runtime.Gosched()
	return deleted
}
