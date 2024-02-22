package repository

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"database/sql"
)

type DataBaseRepo interface {
	Connection() *sql.DB
	AllUsers() ([]*data.User, error)
	GetUser(id int) (*data.User, error)
	GetUserByEmail(email string) (*data.User, error)
	UpdateUser(u data.User) error
	DeleteUser(id int) error
	InsertUser(user data.User) (int, error)
	ResetPassword(id int, password string) error
	InsertUserImage(i data.UserImage) (int, error)
}
