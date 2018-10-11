package model

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

type RawTodo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Todo struct {
	ID        int64     `json:"-" gorm:"type:bigint;primary_key"`
	Title     string    `json:"title" gorm:"type:varchar(80)"`
	Completed bool      `json:"completed" gorm:"type:bool"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Group     Group     `json:"group,omitempty"`
	Owner     User      `json:"owner,omitempty"`
	Users     []*User   `json:"-" gorm:"many2many:user_todos"`
}

type Buddy struct {
	Buddy int64 `json:"buddy"`
}

type TodoServiceImpl interface {
	Get(id int64) (Todo, error)
	GetAll() ([]Todo, error)
	New(todo RawTodo, user User) (Todo, error)
	AddBuddy(todo *Todo, buddy *User) (*Todo, error)
}

type TodoService struct {
	DB        *gorm.DB
	Snowflake *snowflake.Node
}

func (s *TodoService) Get(id int64) (Todo, error) {
	todo := Todo{}
	res := s.DB.First(&todo, id)

	if res.Error != nil || res.RecordNotFound() {
		return todo, res.Error
	}
	return todo, nil
}

func (s *TodoService) GetAll() ([]Todo, error) {
	todos := []Todo{}
	res := s.DB.Find(&todos)

	if res.Error != nil {
		return todos, res.Error
	}
	return todos, nil
}

func (s *TodoService) New(todo RawTodo, user User) (Todo, error) {
	built := s.decorateNewTodo(todo, user, s.Snowflake)

	fmt.Println(built)

	create := s.DB.Create(&built)

	if create.Error != nil {
		fmt.Println(create.Error)
		return built, create.Error
	}

	res := s.DB.Save(&built)

	if res.Error != nil {
		return built, res.Error
	}
	return built, nil
}

func (s *TodoService) decorateNewTodo(i RawTodo, user User, snowflake *snowflake.Node) Todo {
	processed := Todo{
		ID:        snowflake.Generate().Int64(),
		Title:     i.Title,
		Owner:     user,
		Completed: i.Completed,
	}

	return processed
}

func (s *TodoService) AddBuddy(todo *Todo, buddy *User) (*Todo, error) {
	assoc := s.DB.Model(&todo).Association("Users").Append(&User{ID: buddy.ID})
	return todo, assoc.Error
}

type Group struct {
	ID        int64     `json:"-" gorm:"type:bigint;primary_key"`
	Title     string    `json:"title" gorm:"type:varchar(40)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Todos     []Todo
}
