package database

type BaseUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRegister struct {
	BaseUser
	Password string `json:"password"`
}

type UserOut struct {
	BaseUser
	ID int
}

type User struct {
	UserOut
	HashedPassword string
}
