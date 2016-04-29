package app

type Service struct {
	Name         string
	Destinations []string
	Endpoint     string
	Job          JobOpt
	Discovery    DiscoveryOpt
	Metrics      MetricsOpt
	Network      NetworkOpt
}

type DiscoveryOpt struct {
	EtcdAddress  string
	UseDiscovery bool
}

type MetricsOpt struct {
	InfluxAddress string
	InfluxDbName  string
	InfluxUser    string
	InfluxPwd     string
}

type NetworkOpt struct {
	Ip      string
	Port    string
	Address string
}

type JobOpt struct {
	Lambda   float64
	Profiles []string
}
