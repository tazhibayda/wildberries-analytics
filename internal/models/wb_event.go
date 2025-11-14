// internal/models/wb_event.go
package models

type WBEvent struct {
	Type       string      `json:"type"` // "sales", "orders", "stocks"
	SupplierID int         `json:"supplier_id"`
	Data       interface{} `json:"data"`
	CreatedAt  string      `json:"created_at"`
	Source     string      `json:"source"` // "wildberries"
}
