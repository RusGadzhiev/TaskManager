package users

import "context"

type User struct {
	UserName     string `bson:"username"`
	Password string `bson:"password"`
}

const (
	UserName      = "username"
	Password      = "password"
)

// хранит структуру юзеров
type UsersRepo interface {
	// добавляет нового юзера, не проверяя существование такого
	Add(ctx context.Context, user *User) error
	// проверяет наличие юзера с именем 
	IsUserExist(ctx context.Context, name string) bool
	// возвращает пароль
	GetPassword(ctx context.Context, name string) (string, error)
}
