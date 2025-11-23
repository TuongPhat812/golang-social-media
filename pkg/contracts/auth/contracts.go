package auth

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // seconds
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // seconds
}

type LogoutRequest struct {
	Token string `json:"token"` // Access token to blacklist
}

type RevokeTokenRequest struct {
	Token string `json:"token"` // Token to revoke (access or refresh)
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

type UpdateProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"userId,omitempty"`
}

type ProfileResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
