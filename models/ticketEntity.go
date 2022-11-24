package models

type Ticket struct {
	ID         uint   `json:"id"`
	Name       string `gorm:"not null" json:"name"`
	Desc       string `gorm:"not null" json:"desc"`
	Allocation uint   `gorm:"not null" json:"allocation"`
}

type TicketPurchase struct {
	Quantity uint `json:"quantity"`
	UserID   uint `json:"user_id"`
}
