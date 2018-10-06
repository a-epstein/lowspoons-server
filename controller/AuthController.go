package controller

import (
	"time"

	"github.com/kataras/iris"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lowspoons-server/model"
)

type AuthController struct {
	Service model.UserServiceImpl
}

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func genToken(name string, session int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	expires := time.Now().Add(time.Hour * 72).Unix()

	claims["name"] = name
	claims["exp"] = expires
	claims["jti"] = session

	// Generate encoded token and send it as response.
	return token.SignedString([]byte("jwj3ofjfewj1j"))
}

func (c *AuthController) GetBy(name string) interface{} {
	user, err := c.Service.GetByName(name)

	if err != nil {
		return ErrorResponse{Error: "User not found"}
	}

	// Generate encoded token and send it as response.
	t, err := genToken(user.Handle, user.SessionID)
	if err != nil {
		return ErrorResponse{Error: "Token creation error"}
	}

	return AuthResponse{Token: t}
}

func (c *AuthController) Post(ctx iris.Context) interface{} {
	ru := model.RawUser{}
	err := ctx.ReadJSON(&ru)

	if err != nil {
		return ErrorResponse{Error: "Incorrect input"}
	}

	user, err := c.Service.New(ru.Handle)

	if err != nil {
		return ErrorResponse{Error: "Error creating user"}
	}

	// Generate encoded token and send it as response.
	t, err := genToken(user.Handle, user.SessionID)
	if err != nil {
		return ErrorResponse{Error: "Token creation error"}
	}

	return AuthResponse{Token: t}
}
