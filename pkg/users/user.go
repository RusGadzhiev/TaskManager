package users

import "context"

type User struct {
	Id       uint64 `bson:"_id"` // зачем это поле
	Login    string `bson:"login"`
	Password string `bson:"password"`
}

type UsersRepo interface {
	Add(ctx context.Context, user *User) error
	Check(ctx context.Context, user *User) (bool, error)
}
