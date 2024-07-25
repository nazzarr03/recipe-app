package models

type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Latency   string `json:"latency"`
	Method    string `json:"method"`
	Path      string `json:"path"`
}
