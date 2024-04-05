package users

import (
	"HW4/internal/config"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrConnectionMongo = errors.New("error of connecting with mongo db")
	ErrPingMongo       = errors.New("error of ping mongo db")
	ErrNoUser          = errors.New("no such user")
	ErrUserExist       = errors.New("user with this login exists")
)

const (
	DBName         = "task_manager"
	CollectionName = "users"
)

// как тут используются контексты

type UsersRepoMongoDB struct {
	DB *mongo.Collection
}

func NewUsersRepoMongoDB(ctx context.Context, cfg *config.MongoDb) (*UsersRepoMongoDB, *mongo.Client) {
	uri := fmt.Sprintf("mongodb://" + cfg.Host + ":" + cfg.Port + "/messenger?directConnection=true")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrConnectionMongo)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingMongo)
	}

	collection := client.Database(DBName).Collection(CollectionName)
	return &UsersRepoMongoDB{DB: collection}, client
}

func (repo *UsersRepoMongoDB) Check(ctx context.Context, user *User) (bool, error) {
	record := &User{}
	res := repo.DB.FindOne(ctx, bson.M{"login": user.Login})
	if res.Err() == mongo.ErrNoDocuments {
		return false, ErrNoUser
	}
	res.Decode(record)
	return true, nil
}

func (repo *UsersRepoMongoDB) Add(ctx context.Context, user *User) error {

	if exist, _ := repo.Check(ctx, user); exist {
		return ErrUserExist
	}

	newItem := bson.M{
		"login":    user.Login,
		"password": user.Password,
	}

	_, err := repo.DB.InsertOne(ctx, newItem)
	if err != nil {
		return fmt.Errorf("insert mongo error: %w", err)
	}
	return nil
}
