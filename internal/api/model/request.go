package model

type Command struct {
	CmdType string `json:"cmd_type" binding:"alpha"`
	Key     string `json:"key"`
	Value   string `json:"value" binding:"omitempty"`
}

type TransactionRequest struct {
	Commands []Command `json:"commands"`
}
