package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrConnectionMongo = errors.New("error of connecting with mongo db")
	ErrPingMongo       = errors.New("error of ping mongo db")
)

const (
	dbName         = "task_manager"
	collectionName = "users"
)

type UsersRepoMongodb struct {
	db *mongo.Collection
}

func NewUsersRepoMongoDB(ctx context.Context, cfg *config.MongoDb) (*UsersRepoMongodb, *mongo.Client) {
	uri := fmt.Sprintf("mongodb://" + cfg.Host + ":" + cfg.Port + "/messenger?directConnection=true")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrConnectionMongo)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error: %s, Description: %s", err, ErrPingMongo)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &UsersRepoMongodb{db: collection}, client
}

func (repo *UsersRepoMongodb) GetUser(ctx context.Context, username string) (*service.User, error) {
	res := repo.db.FindOne(ctx, bson.M{service.UserName: username})
	if res.Err() == mongo.ErrNoDocuments {
		return nil, service.ErrNoUser
	} else if res.Err() != nil {
		return nil, fmt.Errorf("get user mongo error: %w", res.Err())
	}
	var user service.User
	err := res.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("get user mongo error (decode): %w", err)
	}
	return &user, nil
}
func (repo *UsersRepoMongodb) AddUser(ctx context.Context, user *service.User) error {
	_, err := repo.db.InsertOne(ctx, *user)
	if err != nil {
		return fmt.Errorf("insert mongo error: %w", err)
	}
	return nil
}
