package controller

import (
	"fmt"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/lowspoons-server/model"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

// TodoController is our TODO app's web controller.
type TodoController struct {
	Service     model.TodoServiceImpl
	UserService model.UserServiceImpl
}

// BeforeActivation called once before the server ran, and before
// the routes and dependencies binded.
// You can bind custom things to the controller, add new methods, add middleware,
// add dependencies to the struct or the method(s) and more.
func (c *TodoController) BeforeActivation(b mvc.BeforeActivation) {
	// this could be binded to a controller's function input argument
	// if any, or struct field if any:
	b.Dependencies().Add(func(ctx iris.Context) (item model.RawTodo) {
		ctx.ReadJSON(&item)
		return
	})
}

// Get handles the GET: /todos route.
func (c *TodoController) Get() []model.Todo {
	todos, err := c.Service.GetAll()

	if err != nil {
		fmt.Println("error retrieving todos")
	}

	return todos
}

func (c *TodoController) GetBy(ctx iris.Context, id int64) interface{} {
	todo, err := c.Service.Get(id)

	if err != nil {
		fmt.Println("error retrieving todo %i", id)
		return iris.StatusNotFound
	}

	return todo
}

// PostItemResponse the response data that will be returned as json
// after a post save action of all todo items.
type PostItemResponse struct {
	Success bool `json:"success"`
}

var emptyResponse = PostItemResponse{Success: false}

// Post handles the POST: /todos route.
func (c *TodoController) Post(ctx iris.Context, rawTodo model.RawTodo) PostItemResponse {
	token := ctx.Values().Get("jwt").(*jwt.Token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ids := claims["jti"].(string)
		id, cerr := strconv.ParseInt(ids, 10, 64)

		if cerr != nil {
			fmt.Println("incorrect id format")
			return emptyResponse
		}

		owner, error := c.UserService.GetBySession(id)

		if error != nil {
			fmt.Println("user not found")
			return emptyResponse
		}

		if _, err := c.Service.New(rawTodo, owner); err != nil {
			return emptyResponse
		}

		return PostItemResponse{Success: true}
	}

	return emptyResponse
}

func (c *TodoController) PutBy(ctx iris.Context, id int64) interface{} {
	bu := model.Buddy{}
	jerr := ctx.ReadJSON(&bu)

	if jerr != nil {
		return ErrorResponse{Error: "Incorrect input"}
	}

	buddy, uerr := c.UserService.Get(bu.Buddy)
	if uerr != nil {
		return emptyResponse
	}

	todo, terr := c.Service.Get(id)
	if terr != nil {
		return emptyResponse
	}

	_, err := c.Service.AddBuddy(&todo, &buddy)

	if err != nil {
		return emptyResponse
	}

	return todo
}
