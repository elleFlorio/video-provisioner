package metric

import (
	"log"
	"time"

	"github.com/elleFlorio/video-provisioner/Godeps/_workspace/src/github.com/influxdb/influxdb/client/v2"
)

type InfluxConfig struct {
	Address  string
	DBname   string
	Username string
	Password string
}

var (
	tags           map[string]string
	execFields     map[string]interface{}
	respTimeFields map[string]interface{}
	reqArrFields   map[string]interface{}
	reqDoneFields  map[string]interface{}
	influx         client.Client
	config         InfluxConfig
	batch          client.BatchPoints
)

func Initialize(serviceName string, serviceAddress string, influxConf InfluxConfig) error {
	log.Println("Initializing metric service...")
	defer log.Println("Done")

	var err error
	tags = map[string]string{
		"name":    serviceName,
		"address": serviceAddress,
	}
	execFields = map[string]interface{}{
		"value": 0.0,
	}
	respTimeFields = map[string]interface{}{
		"value": 0.0,
	}
	reqArrFields = map[string]interface{}{
		"value": 0.0,
	}
	reqDoneFields = map[string]interface{}{
		"value": 0.0,
	}
	config = influxConf

	influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return err
	}

	batch, err = client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.DBname,
		Precision: "ms",
	})
	if err != nil {
		return err
	}

	log.Println("Connected to metric service at address " + serviceAddress)

	return nil
}

func SendExecutionTime(execTime float64) error {
	execFields["value"] = execTime
	point, err := client.NewPoint("execution_time", tags, execFields, time.Now())
	if err != nil {
		return err
	}

	batch.AddPoint(point)
	influx.Write(batch)
	return nil
}

func SendResponseTime(respTime float64) error {
	respTimeFields["value"] = respTime
	point, err := client.NewPoint("response_time", tags, respTimeFields, time.Now())
	if err != nil {
		return err
	}

	batch.AddPoint(point)
	influx.Write(batch)
	return nil
}

func SendRequestsArrived(rpm int) error {
	reqArrFields["value"] = rpm
	point, err := client.NewPoint("req_arr", tags, reqArrFields, time.Now())
	if err != nil {
		return err
	}

	batch.AddPoint(point)
	influx.Write(batch)
	return nil
}

func SendRequestsDone(rpm int) error {
	reqDoneFields["value"] = rpm
	point, err := client.NewPoint("req_done", tags, reqDoneFields, time.Now())
	if err != nil {
		return err
	}

	batch.AddPoint(point)
	influx.Write(batch)
	return nil
}
