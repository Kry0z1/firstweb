package database

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateUserHandler(c *gin.Context) {
	var user UserRegister

	err := json.NewDecoder(c.Request.Body).Decode(&user)

	if err != nil || len(user.Username) == 0 || len(user.Password) == 0 {
		c.String(http.StatusBadRequest, "invalid username and/or password")
		return
	}

	_, err = GetUserByUsername(user.Username)
	if err == nil {
		c.String(http.StatusFound, "user with such username already exists")
		return
	}

	storeUser := User{UserOut{user.BaseUser, 0}, ""}

	err = CreateUserWithHashingPassword(&storeUser, user.Password)

	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Failed to load user to database")
		return
	}

	c.JSON(http.StatusOK, storeUser.UserOut)
}
