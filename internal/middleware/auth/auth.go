package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	crud "github.com/Kry0z1/firstweb/internal/database"
	hasher "github.com/Kry0z1/firstweb/internal/passwordhasher"
)

type token struct {
	AccessToken string
	TokenType   string
}

type tokenData struct {
	Username string
}

type userForAuth struct {
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
}

func authUser(username, password string) (*crud.User, bool) {
	user, err := crud.GetUserByUsername(username)

	if err != nil || !hasher.VerifyPassword(password, user.HashedPassword) {
		return &crud.User{}, false
	}

	return user, true
}

func ContextUser(c context.Context) (*crud.User, bool) {
	user, ok := c.Value("contextUser").(*crud.User)
	return user, ok && user != nil
}

func CheckAuth(t Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("WWW-Authenticate", "Bearer")
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.Next()
			return
		}

		splitted := strings.Split(authorizationHeader, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			c.String(http.StatusUnauthorized, "invalid Authorization header format")
			c.Abort()
			return
		}

		token := splitted[1]

		user, err := t.CheckToken(c, token)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				c.String(http.StatusUnauthorized, "invalid token")
			} else {
				c.String(http.StatusUnauthorized, fmt.Sprintf("internal server error: %v", err.Error()))
			}
			c.Abort()
			return
		}

		c.Set("contextUser", user)
		c.Next()
	}
}

func LoginForToken(t Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("WWW-Authenticate", "Bearer")

		var user userForAuth
		err := c.BindJSON(&user)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Failed to parse username and password: %v", err.Error()))
			return
		}

		_, ok := authUser(user.Username, user.Password)
		if !ok {
			c.String(http.StatusUnauthorized, "Wrong password or username")
			return
		}

		authToken, err := t.CreateToken(map[string]string{"sub": user.Username}, 0)
		if err != nil {
			c.String(http.StatusUnauthorized, fmt.Sprintf("Failed to create token: %v", err.Error()))
			return
		}
		c.JSON(http.StatusOK, token{
			AccessToken: authToken,
			TokenType:   "Bearer",
		})
	}

}
