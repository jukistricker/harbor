package agent

import (
	"harbor/internal/module/agent/handler"
	"harbor/internal/module/metric"
	"net/http"
)

type Module struct {
	WSHandler      http.HandlerFunc
	ControlHandler http.HandlerFunc
}

func NewModule(metricMod *metric.Module) *Module {
	return &Module{
		WSHandler:      handler.WSHandler(metricMod),
		ControlHandler: handler.ControlHandler,
	}
}