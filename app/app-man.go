package app

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/elleFlorio/video-provisioner/discovery"
	"github.com/elleFlorio/video-provisioner/job"
	"github.com/elleFlorio/video-provisioner/logger"
	"github.com/elleFlorio/video-provisioner/metric"
	"github.com/elleFlorio/video-provisioner/network"
	"github.com/elleFlorio/video-provisioner/request"
	"github.com/elleFlorio/video-provisioner/signal"
)

const (
	messagePath  = "/message"
	responsePath = "/response"
)

var (
	service Service
	ch_req  chan request.Request
	ch_stop chan struct{}

	ErrNoDestinations = errors.New("No destinations available")
)

func init() {
	ch_req = make(chan request.Request)
	ch_stop = make(chan struct{})
}

func CreateService(params Service) {
	service = params
}

func StartService() {
	if service.Name == "" {
		log.Fatal("Service cannot stast because has not been created")
	}

	setNetworkAddress()
	startSigsMonitor(ch_stop)
	startJobsManager(ch_req)
	initializeMetricService()
	startDiscovery()

	http.HandleFunc(responsePath, readResponse)
	http.HandleFunc(messagePath, readMessage)

	log.Println("Waiting for requests...")
	log.Fatal(http.ListenAndServe(service.Network.Port, nil))
}

func setNetworkAddress() {
	service.Network.Address = network.GenerateAddress(service.Network.Ip, service.Network.Port)
}

func createLogger(name string) {
	logger.New(name)
}

func startSigsMonitor(ch_stop chan struct{}) {
	go signal.MonitorSignals(ch_stop, service.Discovery.UseDiscovery)
}

func startJobsManager(ch_req chan request.Request) {
	go job.ManageJobs(ch_req, service.Destinations)
}

func initializeMetricService() {
	config := metric.InfluxConfig{
		service.Metrics.InfluxAddress,
		service.Metrics.InfluxDbName,
		service.Metrics.InfluxUser,
		service.Metrics.InfluxPwd,
	}
	err := metric.Initialize(service.Name, service.Network.Ip, config)
	if err != nil {
		log.Fatalf("Error: %s; failded to initialize metric service", err.Error())
	}
}

func startDiscovery() {
	var err error
	err = discovery.InitializeEtcd(service.Discovery.EtcdAddress)
	if err != nil {
		log.Fatalln("Cannot connect to etcd server at ", service.Discovery.EtcdAddress)
	}
	log.Println("Connected to etcd server at ", service.Discovery.EtcdAddress)

	if service.Discovery.UseDiscovery {
		err = discovery.RegisterToEtcd(service.Name, service.Network.Address)
		if err != nil {
			log.Fatalln("Cannot register to etcd server", service.Discovery.EtcdAddress)
		}

		log.Println("Registered to etcd server")

		go discovery.KeepAlive(ch_stop)
	}
}

func readMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Create the request
	nDestination := len(service.Destinations)
	req, err := request.CreateReq(r, nDestination)
	if err != nil {
		log.Println("Cannot read message")
		w.WriteHeader(422)
		return
	}

	// Start work
	ch_req <- req

	w.WriteHeader(http.StatusCreated)
}

func readResponse(w http.ResponseWriter, r *http.Request) {
	var err error
	var respTimeMs float64

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	message, err := network.ReadMessage(r)
	if err != nil {
		log.Println("Cannot read message")
		w.WriteHeader(422)
		return
	}

	reqId := message.Args
	if req, ok := request.GetRequest(reqId); ok {
		respTimeMs = time.Since(req.Start).Seconds() * 1000
		complete := request.UpdateRequestInHistory(reqId)
		if complete {
			network.RespondeToRequest(req.From, req.ID, message.Body)
		}
	} else {
		log.Println("Cannot find request ID in history")
		w.WriteHeader(422)
		return
	}
	if message.Body == "done" {
		logger.LogResponseTime(respTimeMs)
		metric.SendResponseTime(respTimeMs)
	} else {
		log.Println("Error: request lost.")
	}

	w.WriteHeader(http.StatusCreated)
}
