package repository
import (
	"database/sql"

	model "harbor/internal/module/monitor/model"
)

type HealthMetricRepo struct {
	db *sql.DB
}

func NewHealthMetricRepo(db *sql.DB) *HealthMetricRepo {
	return &HealthMetricRepo{db: db}
}

func (r *HealthMetricRepo) SaveBatch(metrics []model.HealthMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO health_metrics (
			endpoint_id, timestamp_ns, status_code, latency_ms, 
			dns_ms, tcp_ms, tls_ms, ttfb_ms, load_ms, res_size_bytes, test_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range metrics {
		_, err := stmt.Exec(
			m.EndpointID, m.TimestampNS, m.StatusCode, m.LatencyMS,
			m.DNSMs, m.TCPMs, m.TLSMs, m.TTFBMs, m.LoadMs, m.ResSizeBytes, m.TestType,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}