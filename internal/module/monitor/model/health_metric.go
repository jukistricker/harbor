package db

// HealthMetric đại diện cho dữ liệu thô thu thập được
type HealthMetric struct {
	EndpointID   int     `json:"endpoint_id"`
	TimestampNS  int64   `json:"timestamp_ns"`
	StatusCode   int     `json:"status_code"`
	LatencyMS    float64 `json:"latency_ms"`
	PrepareMs    float64 `json:"prepare_ms"`
	DNSMs        float64 `json:"dns_ms"`
	TCPMs        float64 `json:"tcp_ms"`
	TLSMs        float64 `json:"tls_ms"`
	TTFBMs       float64 `json:"ttfb_ms"`
	LoadMs       float64 `json:"load_ms"`
	ResSizeBytes int     `json:"res_size_bytes"`
	TestType     int     `json:"test_type"` // 1: Health, 2: Stress
}