package payload

type SignupPayload struct {
	UserName string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=128"`
}

type EmailPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordPayload struct {
	SignupToken string `json:"signupToken" validate:"required"`
	Password    string `json:"password" validate:"required,min=6,max=128"`
}

type OTPPayload struct {
	SignupToken string `json:"signupToken" validate:"required"`
	OTP         int    `json:"otp" validate:"required,numeric"`
}
