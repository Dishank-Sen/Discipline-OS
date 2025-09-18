package types

import "time"

type User struct {
	ID         string                  `json:"id"`
	UserName   string                  `json:"username"`
	Email      string                  `json:"email"`
	Password   string                  `json:"password"`
	Subscribed bool                    `json:"subscribed"`
	Personal   UserPersonalInformation `json:"personal"`
	Device     UserDeviceInformation   `json:"device"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type TempUser struct {
	ID          string    `json:"id"`
	SignupToken string    `json:"signupToken"`
	UserName    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	OTP         int       `json:"otp"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserPersonalInformation struct {
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

type UserDeviceInformation struct {
	Theme              string `json:"theme"`
	NotificationStatus bool   `json:"notificationStatus"`
	TimeZone           string `json:"timezone"`
}
