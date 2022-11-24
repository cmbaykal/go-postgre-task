package models

type Ticket struct {
	ID         int    `json:"id"`
	Name       string `gorm:"not null" json:"name"`
	Desc       string `gorm:"not null" json:"desc"`
	Allocation int    `gorm:"not null" json:"allocation"`
}

type TicketPurchase struct {
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
}
