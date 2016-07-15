package logger

import (
	"log"
	"strconv"
)

var name string

func New(serviceName string) {
	name = serviceName
}

func LogExecutionTime(execTime float64) {
	log.Println("gru:" + name + ":" + "execution_time:" + strconv.FormatFloat(execTime, 'f', 2, 64) + ":ms")
}

func LogResponseTime(respTime float64) {
	log.Println("gru:" + name + ":" + "response_time" + ":" + strconv.FormatFloat(respTime, 'f', 2, 64) + ":ms")
}

func LogRequestsArrivedPerMinute(rpm int) {
	log.Println("gru:" + name + ":" + "rpm_arr" + ":" + strconv.Itoa(rpm) + ":short")
}
