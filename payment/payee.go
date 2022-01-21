package main

type Account struct {
	ID            string `json:"ID"`
	Collection    string `json:"collection"`
	SortCode      string `json:"sortcode"`
	AccountNumber string `json:"accnumber"`
	IsActive      bool   `json:"isActive"`
	WarningID     string `json:"warningID"`
}
