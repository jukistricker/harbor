package health

import (
	"database/sql"
	"harbor/internal/module/health/repository"
	// "harbor/internal/module/health/service"
)

type Module struct {
	Repo    *repository.HealthSnapshotRepo
	// Service *service.HealthSnapshotService
	// Controller *controller.HealthSnapshotController
}

func New(db *sql.DB) *Module {
	repo := repository.NewHealthSnapshotRepo(db)
	// svc := service.NewHealthSnapshotService(repo)
	
	return &Module{
		Repo:    repo,
		// Service: svc,
	}
}