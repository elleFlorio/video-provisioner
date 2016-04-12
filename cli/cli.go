package cli

import (
	"fmt"
	"os"

	"github.com/elleFlorio/video-provisioner/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func Run() {
	app := cli.NewApp()
	app.Name = "testApp"
	app.Usage = "Test application"

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Start the a service",
			Action: start,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "etcdserver, e",
					Usage:  fmt.Sprintf("url of etcd server"),
					EnvVar: "ETCD_ADDR",
				},
				cli.StringFlag{
					Name:   "ipaddress, a",
					Value:  "",
					Usage:  fmt.Sprintf("Ip address of the host"),
					EnvVar: "HostIP",
				},
				cli.StringFlag{
					Name:   "influxdb, m",
					Usage:  fmt.Sprintf("url of influxdb"),
					EnvVar: "INFLUX_ADDR",
				},
				cli.StringFlag{
					Name:   "db-user, dbu",
					Usage:  fmt.Sprintf("influxdb user username"),
					EnvVar: "INFLUX_USER",
				},
				cli.StringFlag{
					Name:   "db-pwd, dbp",
					Usage:  fmt.Sprintf("influxdb user password"),
					EnvVar: "INFLUX_PWD",
				},
				cli.StringFlag{
					Name:  "db-name, db",
					Value: "video-provisioner",
					Usage: fmt.Sprintf("influxdb database name. Default is 'video-provisioner'"),
				},
				cli.StringFlag{
					Name:  "port, p",
					Value: "",
					Usage: fmt.Sprintf("port of the service"),
				},
				cli.Float64Flag{
					Name:  "lambda, l",
					Value: 0.0,
					Usage: fmt.Sprintf("lambda in milliseconds. Lambda is used to" +
						"compute the workload. Default is 0.0"),
				},
				cli.StringSliceFlag{
					Name:  "profiles, lp",
					Value: &cli.StringSlice{},
					Usage: fmt.Sprintf("load profiles to compute the load. " +
						"Each load profiles should be in the format filename:probability"),
				},
				cli.BoolFlag{
					Name:  "discovery, ds",
					Usage: fmt.Sprintf("Register to discovery service"),
				},
				cli.StringSliceFlag{
					Name:  "destinations, d",
					Value: &cli.StringSlice{},
					Usage: fmt.Sprintf("destination of request messages. " +
						"Each destination should be in the format service:probability"),
				},
			},
		},
	}

	app.Run(os.Args)
}
