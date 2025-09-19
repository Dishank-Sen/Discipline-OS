package types

import "time"

type User struct {
	ID         string                  `json:"id" validate:"required"`
	UserName   string                  `json:"username" validate:"required,min=3,max=50"`
	Email      string                  `json:"email" validate:"required,email"`
	Password   string                  `json:"password" validate:"required,min=6,max=128"`
	Subscribed bool                    `json:"subscribed"`
	Personal   UserPersonalInformation `json:"personal" validate:"required,dive"`
	Device     UserDeviceInformation   `json:"device" validate:"required,dive"`
	CreatedAt  time.Time               `json:"createdAt" validate:"required" default:"now"`
}

type TempUser struct {
	ID          string    `json:"id" validate:"required"`
	SignupToken string    `json:"signupToken" validate:"required"`
	UserName    string    `json:"username" validate:"required,min=3,max=50"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=6,max=128"`
	OTP         int       `json:"otp" validate:"required,numeric"`
	UpdatedAt   time.Time `json:"updatedAt" validate:"required" default:"now"`
}

type UserPersonalInformation struct {
	Age    int    `json:"age" validate:"required,min=0,max=120"`
	Gender string `json:"gender" validate:"required,oneof=male female other"`
	Height int    `json:"height" validate:"omitempty,min=0,max=300"`
	Weight int    `json:"weight" validate:"omitempty,min=0,max=500"`
}

type UserDeviceInformation struct {
	Theme              string `json:"theme" validate:"required,oneof=light dark"`
	NotificationStatus bool   `json:"notificationStatus"`
	TimeZone           string `json:"timezone" validate:"required"`
}
