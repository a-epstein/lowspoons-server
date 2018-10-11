package model

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jinzhu/gorm"
)

type RawUser struct {
	Handle string `json:"handle"`
}

type User struct {
	ID        int64     `json:"-" gorm:"type:bigint;primary_key"`
	SessionID int64     `json:"session,omitempty" gorm:"type:bigint"`
	Handle    string    `json:"handle" gorm:"type:varchar(28);unique_index"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Todos     []Todo    `json:"-" gorm:"many2many:user_todos"`
}

type UserServiceImpl interface {
	New(handle string) (User, error)
	GetByName(handle string) (User, error)
	Get(id int64) (User, error)
	GetBySession(id int64) (User, error)
}

type UserService struct {
	DB        *gorm.DB
	Snowflake *snowflake.Node
}

func (u *UserService) newUserBuilder(handle string) User {
	return User{
		ID:        u.Snowflake.Generate().Int64(),
		SessionID: u.Snowflake.Generate().Int64(),
		Handle:    handle,
	}
}

func (u *UserService) Get(id int64) (User, error) {
	user := User{}
	res := u.DB.First(&user, id)

	if res.Error != nil || res.RecordNotFound() {
		return user, res.Error
	}
	return user, nil
}

func (u *UserService) GetByName(name string) (User, error) {
	user := User{}
	res := u.DB.Where(&User{Handle: name}).First(&user)

	if res.Error != nil || res.RecordNotFound() {
		return user, res.Error
	}
	return user, nil
}

func (u *UserService) GetBySession(id int64) (User, error) {
	user := User{}
	res := u.DB.Where(&User{SessionID: id}).First(&user)

	if res.Error != nil || res.RecordNotFound() {
		return user, res.Error
	}
	return user, nil
}

func (u *UserService) New(handle string) (User, error) {

	existing := User{}
	u.DB.Where("handle = ?", handle).First(&existing)

	if existing.ID > 0 {
		return existing, nil
	}

	user := u.newUserBuilder(handle)
	u.DB.Create(&user)
	res := u.DB.Save(&user)

	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}
