package api

type Task struct {
	RunnerId string `json:"runner_id"`
	Content  string `json:"content"`
	Lang     string `json:"lang"`
	Hash     uint32 `json:"hash"`
	Result   string `json:"result"`
}
