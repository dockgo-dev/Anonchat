package models

type (
	DefaultResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	TokenData struct {
		AcessToken   string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	UserData struct {
		UserID int64  `json:"user_id"`
		Login  string `json:"login"`
		Email  string `json:"email"`
	}

	RegisterRequest struct {
		Login    string `json:"login"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	TokenRequest struct {
		Token string `json:"token"`
	}
)
