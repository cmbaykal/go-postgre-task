package models

// swagger:model
type Ticket struct {
	ID         int    `json:"id"`
	Name       string `gorm:"not null" json:"name"`
	Desc       string `gorm:"not null" json:"desc"`
	Allocation int    `gorm:"not null" json:"allocation"`
}

// swagger:model
type TicketPurchase struct {
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
}
