package service

type ChargeJSON struct {
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receipt_email"`
}
