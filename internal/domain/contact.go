package domain

import "time"

type Contact struct {
	UserID    string    `json:"userId"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type ContactResponse struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
}

type SyncContactsRequest struct {
	Contacts []ContactInput `json:"contacts" binding:"required"`
}

type ContactInput struct {
	Phone string `json:"phone" binding:"required"`
	Name  string `json:"name" binding:"required"`
}
