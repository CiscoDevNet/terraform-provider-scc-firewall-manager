package transaction

import (
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/publicapi/transaction/transactiontype"
)

type Type struct {
	TransactionUid  string                 `json:"transactionUid"`
	TenantUid       string                 `json:"tenantUid"`
	EntityUid       string                 `json:"entityUid"`
	EntityUrl       string                 `json:"entityUrl"`
	PollingUrl      string                 `json:"transactionPollingUrl"`
	SubmissionTime  string                 `json:"submissionTime"`
	LastUpdatedTime string                 `json:"lastUpdatedTime"`
	Type            transactiontype.Type   `json:"transactionType"`
	Status          transactionstatus.Type `json:"cdoTransactionStatus"`
	Details         map[string]string      `json:"transactionDetails"`
	ErrorMessage    string                 `json:"errorMessage"`
	ErrorDetails    map[string]string      `json:"errorDetails"`
}
