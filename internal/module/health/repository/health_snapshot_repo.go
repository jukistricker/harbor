package repository
import (
	"database/sql"

	model "harbor/internal/module/health/model"
)

type HealthSnapshotRepo struct {
	db *sql.DB
}

func NewHealthSnapshotRepo(db *sql.DB) *HealthSnapshotRepo {
	return &HealthSnapshotRepo{db: db}
}

func (r *HealthSnapshotRepo) GetDashboardData() ([]model.HealthSnapshot, error) {
	query := `
		SELECT 
			t.id, t.name, t.base_url,
			s.last_status, s.last_latency, s.cpu_usage, s.ram_usage_mb, 
			s.health_score, s.bottleneck, s.last_check_at
		FROM monitor_targets t
		LEFT JOIN health_snapshots s ON t.id = s.target_id
		WHERE t.is_active = 1`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.HealthSnapshot
	for rows.Next() {
		var d model.HealthSnapshot
		var lastCheck int64
		err := rows.Scan(
			&d.ID, &d.Name, &d.BaseURL,
			&d.Status, &d.Latency, &d.CPU, &d.RAM,
			&d.HealthScore, &d.Bottleneck, &lastCheck,
		)
		if err == nil {
			d.IsUp = (d.Status >= 200 && d.Status < 400)
			// Bạn có thể thêm logic tính "2 seconds ago" ở đây hoặc để FE làm
			list = append(list, d)
		}
	}
	return list, nil
}