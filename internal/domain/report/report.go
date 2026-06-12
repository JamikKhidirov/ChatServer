package reportdomain

type MessageReport struct {
	ID          string `json:"id"`
	MessageID   string `json:"message_id"`
	ReporterID  string `json:"reporter_id"`
	Reason      string `json:"reason"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ResolvedBy  string `json:"resolved_by,omitempty"`
	ResolvedAt  string `json:"resolved_at,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type CreateReportRequest struct {
	MessageID   string `json:"message_id" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	Description string `json:"description"`
}

type ResolveReportRequest struct {
	Status string `json:"status" binding:"required,oneof=pending resolved dismissed"`
}
