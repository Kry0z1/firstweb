package database

import hasher "github.com/Kry0z1/firstweb/internal/passwordhasher"

func CreateUser(user *User) error {
	err := db.Create(user).Error
	return err
}

func CreateUserWithHashingPassword(user *User, password string) error {
	hashed, err := hasher.GetPasswordHash(password)

	if err != nil {
		return err
	}

	user.HashedPassword = hashed
	return CreateUser(user)
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.First(&user, "username=?", username).Error
	return &user, err
}

func DeleteUserByUsername(username string) error {
	user, err := GetUserByUsername(username)
	if err != nil {
		return err
	}

	err = db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}
