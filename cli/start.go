package cli

import (
	"log"
	"strconv"

	"github.com/elleFlorio/video-provisioner/Godeps/_workspace/src/github.com/codegangsta/cli"

	"github.com/elleFlorio/video-provisioner/app"
	"github.com/elleFlorio/video-provisioner/network"
)

func start(c *cli.Context) {
	if !c.Args().Present() {
		log.Fatalln("Cannot start service: service name is missing")
	}

	name := c.Args().First()
	destinations := c.StringSlice("destinations")

	lambda := c.Float64("lambda")
	profiles := c.StringSlice("profiles")
	job := app.JobOpt{
		Lambda:   lambda,
		Profiles: profiles,
	}

	etcdAddress := c.String("etcdserver")
	useDiscovery := c.Bool("discovery")
	discovery := app.DiscoveryOpt{
		EtcdAddress:  etcdAddress,
		UseDiscovery: useDiscovery,
	}

	influxAddress := c.String("influxdb")
	influxDB := c.String("db-name")
	influxUser := c.String("db-user")
	influxPwd := c.String("db-pwd")
	metrics := app.MetricsOpt{
		InfluxAddress: influxAddress,
		InfluxDbName:  influxDB,
		InfluxUser:    influxUser,
		InfluxPwd:     influxPwd,
	}

	var ip string
	if ip = c.String("ipaddress"); ip == "" {
		ip = network.GetHostIp()
	}
	var port string
	if port = c.String("port"); port == "" {
		p := network.GetPort()
		port = strconv.Itoa(p)
	}
	port = ":" + port
	network := app.NetworkOpt{
		Ip:   ip,
		Port: port,
	}

	service := app.Service{
		Name:         name,
		Destinations: destinations,
		Job:          job,
		Discovery:    discovery,
		Metrics:      metrics,
		Network:      network,
	}

	app.CreateService(service)

	app.StartService()
}
