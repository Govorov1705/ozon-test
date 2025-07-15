package dtos

type AuthRequest struct {
	Username string `validate:"min=6,max=100"`
	Password string `validate:"min=6"`
}
