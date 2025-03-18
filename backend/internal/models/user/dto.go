package user

type AuthorizeInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthorizeOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignUpInput struct {
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Name       string  `json:"name"`
	IsTeacher  int     `json:"is_teacher"`
	InviteCode *string `json:"invite_code,omitempty"`
}

type SignUpOutput struct {
	Success bool `json:"success"`
}
