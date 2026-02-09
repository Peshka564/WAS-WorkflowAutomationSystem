package dto

type Template struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	EmailTo string `json:"email_to" validate:"required"`
}
