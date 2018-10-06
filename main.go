package main

import (
	"fmt"

	"github.com/lowspoons-server/model"

	"github.com/kataras/iris"
	"github.com/lowspoons-server/controller"
	"github.com/lowspoons-server/service"

	"github.com/bwmarrin/snowflake"
	"github.com/kataras/iris/mvc"

	jwt "github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

func main() {
	app := iris.New()

	// serve our app in public, public folder
	// contains the client-side vue.js application,
	// no need for any server-side template here,
	// actually if you're going to just use vue without any
	// back-end services, you can just stop afer this line and start the server.
	app.StaticWeb("/", "./public")

	snowflake, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := new(service.PGService)
	_, dbErr := db.Connect()

	if dbErr != nil {
		fmt.Println("DB connection failed")
		return
	}

	userservice := &model.UserService{DB: db.DB, Snowflake: snowflake}
	todoservice := &model.TodoService{DB: db.DB, Snowflake: snowflake}

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("jwj3ofjfewj1j"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	mvc.Configure(app.Party("/api/todo"), func(app *mvc.Application) {
		app.Router.Use(jwtHandler.Serve)
		app.Register(todoservice)
		app.Register(userservice)
		app.Handle(new(controller.TodoController))
	})

	mvc.Configure(app.Party("/api/auth"), func(app *mvc.Application) {
		app.Register(userservice)
		app.Handle(new(controller.AuthController))
	})

	// start the web server at http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutVersionChecker, iris.WithOptimizations)
}
