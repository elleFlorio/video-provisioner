package network

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
)

var (
	destinations map[string]float64
	rnd          *rand.Rand
)

func init() {
	destinations = make(map[string]float64)
	source := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(source)
}

func ReadDestinations(destString []string) {
	if len(destString) == 0 {
		return
	}

	log.Println("Reading destinations...")
	defer log.Println("Done")

	var prob float64
	var err error
	probSum := 0.0

	for _, dest := range destString {
		destProb := strings.Split(dest, ":")
		if len(destProb) < 2 {
			prob = 1.0
		} else {
			prob, err = strconv.ParseFloat(destProb[1], 64)
			if err != nil {
				log.Println("Error parsing destination probability. Set to 0.0")
				prob = 0.0
			}
		}

		destinations[destProb[0]] = prob

		probSum += prob

		log.Printf("Added destination %s with probability %f\n", destProb[0], prob)
	}

	if probSum != 1.0 {
		log.Fatalln("Error: the sum of destinations probabilities should be 1.0")
	}
}

func GetDestinationsNumber() int {
	return len(destinations)
}

func SendMessageToSpecificService(requestID string, service string) error {
	instances, err := discovery.GetAvailableInstances(service)
	if err != nil {
		log.Println("Cannot dispatch message to service ", service)
		return err
	}
	instance := getInstance(instances)
	sendReqToDest(requestID, instance)
	return nil
}

func SendMessageToDestination(requestID string) error {
	destination := getDestination()
	instances, err := discovery.GetAvailableInstances(destination)
	if err != nil {
		log.Println("Cannot dispatch message to service ", destination)
		return errors.New("Cannot dispatch message to service " + destination)
	}
	instance := getInstance(instances)
	sendReqToDest(requestID, instance)

	return nil
}

func getDestination() string {
	if len(destinations) == 1 {
		for dest, _ := range destinations {
			return dest
		}
	}

	p := rnd.Float64()
	probSum := 0.0
	for dest, prob := range destinations {
		probSum += prob
		if p <= probSum {
			return dest
		}
	}

	log.Println("Error: unable to get a destination")
	return ""
}

func getInstance(instances []string) string {
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
