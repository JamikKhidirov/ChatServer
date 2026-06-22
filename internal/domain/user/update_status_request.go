package userdomain

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"max=100"`
}
