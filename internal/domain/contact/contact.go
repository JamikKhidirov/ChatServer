package contactdomain

import "time"

type Contact struct {
	UserID    string    `json:"userId"`
	Phone     string    `json:"phone"`
	Name      string    `json:"name"`
	PhotoURL  string    `json:"photoUrl,omitempty"`
	UserIDRef string    `json:"userIdRef,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type ContactResponse struct {
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	PhotoURL  string `json:"photoUrl,omitempty"`
	UserIDRef string `json:"userIdRef,omitempty"`
}

type SyncContactsRequest struct {
	Contacts []ContactInput `json:"contacts" binding:"required"`
}

type ContactInput struct {
	Phone string `json:"phone" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

type UpdateContactPhotoRequest struct {
	Phone    string `json:"phone" binding:"required"`
	PhotoURL string `json:"photoUrl" binding:"required"`
}
