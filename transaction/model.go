package transaction

// RunningBalance is a struct to hold balance details for a
// set of transactions for an entity
type RunningBalance struct {
	TransactionID         string `json:"transaction_id,omitempty"`
	CurrentRunningBalance int64  `json:"current_running_balance"`
	TotalCredit           int64  `json:"total_credit"`
	TotalDebit            int64  `json:"total_debit"`
}
