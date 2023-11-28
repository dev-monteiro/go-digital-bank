package business

type CurrInvoiceResp struct {
	StatusLabel    string `json:"status"`
	Amount         string `json:"amount"`
	FmtClosingDate string `json:"closingDate"`
}
