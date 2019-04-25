package adminapi

import (
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	ServiceOnline  = "online"
	ServiceOffline = "offline"
)

//return status of services
type IHealthCheckMonitor interface {
	CheckWeb() string
	CheckApi() string
	CheckEngine() string
	CheckLauncher() string
	CheckWatcher() string
	CheckWebSocket() string
}

type HealthCheckService struct {
	client utils.IHttpClient

	options *HealthCheckOptions
}

type HealthCheckOptions struct {
	webUrl       string
	apiUrl       string
	engineUrl    string
	launcherUrl  string
	watcherUrl   string
	websocketUrl string
}

func NewHealthCheckService(options *HealthCheckOptions) IHealthCheckMonitor {
	if options == nil {
		options = &HealthCheckOptions{
			webUrl:       os.Getenv("WEB_HEALTH_CHECK_URL"),
			apiUrl:       os.Getenv("API_HEALTH_CHECK_URL"),
			engineUrl:    os.Getenv("ENGINE_HEALTH_CHECK_URL"),
			launcherUrl:  os.Getenv("LAUNCHER_HEALTH_CHECK_URL"),
			watcherUrl:   os.Getenv("WATCHER_HEALTH_CHECK_URL"),
			websocketUrl: os.Getenv("WEBSOCKET_HEALTH_CHECK_URL"),
		}
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 500 * time.Millisecond,
		}).DialContext,
		TLSHandshakeTimeout: 1000 * time.Millisecond,
	}
	return &HealthCheckService{
		client:  utils.NewHttpClient(transport),
		options: options,
	}
}

func (h *HealthCheckService) CheckWeb() string {
	_, code, _ := h.client.Get(h.options.webUrl, nil, nil, nil)
	return ToStatus(code)
}

func (h *HealthCheckService) CheckApi() string {
	_, code, _ := h.client.Get(h.options.apiUrl, nil, nil, nil)
	return ToStatus(code)
}

func (h *HealthCheckService) CheckEngine() string {
	_, code, _ := h.client.Get(h.options.engineUrl, nil, nil, nil)
	return ToStatus(code)
}

func (h *HealthCheckService) CheckLauncher() string {
	_, code, _ := h.client.Get(h.options.launcherUrl, nil, nil, nil)
	return ToStatus(code)
}

func (h *HealthCheckService) CheckWatcher() string {
	_, code, _ := h.client.Get(h.options.watcherUrl, nil, nil, nil)
	return ToStatus(code)
}

func (h *HealthCheckService) CheckWebSocket() string {
	_, code, _ := h.client.Get(h.options.websocketUrl, nil, nil, nil)
	return ToStatus(code)
}

func ToStatus(httpCode int) string {
	if http.StatusOK == httpCode {
		return ServiceOnline
	}

	return ServiceOffline
}
