package models

type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	IsCompletedSetup bool   `json:"is_completed_setup"`
	SessionToken     string `json:"session_token"`
	SessionExpired   string `json:"session_expired"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdate struct {
	Id               int     `json:"id" validate:"required"`
	Username         string  `json:"username,omitempty" validate:"min=5,max=32,alphanum"`
	Password         string  `json:"password,omitempty" validate:"min=6,max=64,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
	PreviousPassword string  `json:"previous_password,omitempty"` // for password update
	IsCompletedSetup *bool   `json:"is_completed_setup,omitempty"`
	SessionToken     *string `json:"session_token,omitempty"`   // pointer for optional/null
	SessionExpired   *string `json:"session_expired,omitempty"` // pointer for optional/null
}
