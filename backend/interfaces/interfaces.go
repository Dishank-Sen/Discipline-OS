package interfaces

import types "github.com/Dishank-Sen/Discipline-OS/types/database"

type UserStore interface{
	GetUserByEmail(email string) (*types.User, error)
}