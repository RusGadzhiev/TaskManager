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
)

const (
	DBName         = "task_manager"
	CollectionName = "users"
)

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

func (repo *UsersRepoMongoDB) IsUserExist(ctx context.Context, name string) bool {
	res := repo.DB.FindOne(ctx, bson.M{UserName: name})
	return res.Err() != mongo.ErrNoDocuments
}

func (repo *UsersRepoMongoDB) GetPassword(ctx context.Context, name string) (string, error) {
	res := repo.DB.FindOne(ctx, bson.M{UserName: name})
	if res.Err() != nil {
		return "", fmt.Errorf("get password mongo error: %w", res.Err())
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		return "", fmt.Errorf("get password mongo error (decode): %w", err)
	}
	return user.Password, nil
}

func (repo *UsersRepoMongoDB) Add(ctx context.Context, user *User) error {
	_, err := repo.DB.InsertOne(ctx, *user)
	if err != nil {
		return fmt.Errorf("insert mongo error: %w", err)
	}
	return nil
}
