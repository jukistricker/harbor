package model

//Health snapshot đại diện cho một mục hiển thị trên dashboard giám sát
type HealthSnapshot struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	BaseURL        string  `json:"base_url"`
	Status         int     `json:"status"`         // 200, 500, etc.
	CPU            float64 `json:"cpu"`            // %
	RAM            float64 `json:"ram"`            // MB
	Latency        float64 `json:"latency"`        // ms
	HealthScore    int     `json:"health_score"`   // 0-100
	Bottleneck     string  `json:"bottleneck"`     // 'CPU', 'Network'...
	IsUp           bool    `json:"is_up"`
	LastCheckHuman string  `json:"last_check"`     // "2 seconds ago"
}