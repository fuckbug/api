package users

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Signup struct {
	Email    string `json:"email" validate:"required,email" example:"me@example.com"`
	Password string `json:"password" validate:"required" example:"1234567890"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"me@example.com"`
	Password string `json:"password" validate:"required" example:"1234567890"`
}

type Token struct {
	AccessToken  string `json:"accessToken" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	RefreshToken string `json:"refreshToken" example:"a08929b5-d4f0-4ceb-9cfe-bb4fc05b030c"`
	ExpiresIn    int64  `json:"expiresIn" example:"300"`
}
