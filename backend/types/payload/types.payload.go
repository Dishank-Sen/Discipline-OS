package payload

type SignupPayload struct{
	UserName   string                  `json:"username"`
	Email      string                  `json:"email"`
	Password   string                  `json:"password"`
}