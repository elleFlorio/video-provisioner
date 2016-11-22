# Video-Provisioner #

This is the use-case application used to test the [Gru](https://github.com/elleFlorio/gru) project, and is based on a modified version of [mu-sim](https://github.com/elleFlorio/mu-sim) project.
The purpose of this project is to simulate the microservices composing an application that provides videos on-demand. The microservices can be configured to send requests between them, creating an execution flow. The load between microservices is balanced using round robin scheduling.

The project requires [etcd](https://github.com/coreos/etcd) and [influxDB](https://github.com/influxdata/influxdb) to run properly.
**etcd** is used for service discovery and to keep track of the requests, allowing the dynamic scaling of the microservices.
**InfluxDB** is used for live monitoring of the response time of the application and the microservices.

the docker image is available as [elleflorio/video](https://hub.docker.com/r/elleflorio/video/).