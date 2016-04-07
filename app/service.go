package app

type Service struct {
	Name         string
	Workload     float64
	Destinations []string
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
