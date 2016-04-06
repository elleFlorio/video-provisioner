package network

import (
	"log"
	"math/rand"

	"github.com/elleFlorio/video-provisioner/discovery"
)

func SendMessageToSpecificService(requestID string, service string) error {
	instances, err := discovery.GetAvailableInstances(service)
	if err != nil {
		log.Println("Cannot dispatch message to service ", service)
		return err
	}
	destination := getDestination(instances)
	sendReqToDest(requestID, destination)
	return nil
}

func SendMessageToDestinations(requestID string, destinations []string) int {
	errCounter := 0

	for _, service := range destinations {
		instances, err := discovery.GetAvailableInstances(service)
		if err != nil {
			log.Println("Cannot dispatch message to service ", service)
			errCounter++
			break
		}
		destination := getDestination(instances)
		sendReqToDest(requestID, destination)
	}

	return errCounter
}

func getDestination(instances []string) string {
	if len(instances) == 1 {
		return instances[0]
	}

	return instances[rand.Intn(len(instances))]
}

func sendReqToDest(reqID string, dest string) {
	go Send(dest, "do", reqID, GetMyAddress(), false)
}

func RespondeToRequest(dest string, reqId string, status string) {
	Send(dest, status, reqId, GetMyAddress(), true)
}
