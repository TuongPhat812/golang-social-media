package dto

type LoginUserCommandRequest struct {
	Email    string
	Password string
}

type LoginUserCommandResponse struct {
	UserID string
	Token  string
}
