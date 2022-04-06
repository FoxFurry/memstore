/*
Package model

Describes HTTP models for requests and responses
*/
package model

// Command model used to unmarshal each command in transaction
type Command struct {
	CmdType string `json:"cmd_type" binding:"alpha"`
	Key     string `json:"key"`
	Value   string `json:"value" binding:"omitempty"`
}

// TransactionRequest describes HTTP request for transaction
type TransactionRequest struct {
	Commands []Command `json:"commands"`
}

// TransactionResponse describes HTTP response for successful transaction
type TransactionResponse struct {
	Results []string `json:"results"`
}
