// Package model model transaction
package model

import "time"

// Transaction model
type Transaction struct {
	ID      string    `json:"id"`
	Account string    `json:"account"`
	Amount  float64   `json:"amount"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
