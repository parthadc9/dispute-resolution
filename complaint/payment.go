package main

type PaymentRequest struct {
	Collection           string `json:"collection"`
	ID                   string `json:"ID"`
	PayerID              string `json:"payerID"`
	BeneficiaryID        string `json:"beneficiaryID"`
	Amount               string `json:"amount"`
	Status               string `json:"status"`
	WarningID            string `json:"warningID"`
	PaymentTransactionID string `json:"paymentTransactionID"`
	PaymentMode          string `json:"paymentmode"`
}
