package model

type TransactionDetail struct {
	AvailableAmount string `json:"availableAmount"`
	IsDormant       string `json:"isDormant"`
}

type AccountInfoResponse struct {
	ResponseCode      string            `json:"responseCode"`
	ResponseDesc      string            `json:"responseDesc"`
	TransactionDetail TransactionDetail `json:"transactionDetail"`
}
