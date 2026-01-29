package monitor

import (
	"database/sql"
	"harbor/internal/module/monitor/repository"
	// "harbor/internal/module/monitor/service"
)

type Module struct {
	Repo    *repository.HealthMetricRepo
	// Service *service.MonitorService
	// Controller *controller.MonitorController
}

func New(db *sql.DB) *Module {
	repo := repository.NewHealthMetricRepo(db)
	// svc := service.NewMonitorService(repo)
	
	return &Module{
		Repo:    repo,
		// Service: svc,
	}
}