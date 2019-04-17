package main

import "github.com/HydroProtocol/hydro-sdk-backend/common"

var QueueService common.IQueue

//to monitor if service alive
var HealthCheckMonitor IHealthCheckMonitor

type IHealthCheckMonitor interface {
	CheckWeb() bool
	CheckApi() bool
	CheckEngine() bool
	CheckLauncher() bool
	CheckWatcher() bool
	CheckWebSocket() bool
}
