package api

type input struct {
	Action   string `json:"action"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Content  string `json:"content"`
	Hash     uint32 `json:"hash"`
	Lang     string `json:"lang"`
	RunnerId string `json:"runner_id"`
	Result   string `json:"result"`
	IsPublic bool   `json:"is_public"`
}
