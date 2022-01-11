package domain

import (
	"fmt"
	"forklift/internal/domain/consts"
	"forklift/internal/domain/entity"
	"forklift/pkg/mongo"
	"github.com/globalsign/mgo/bson"
)

type UserRepository interface {
	Upsert(user *entity.User) error
}

type userRepository struct {
	mongoClient  mongo.Client
	databaseName string
}

func (r *userRepository) Upsert(user *entity.User) error {

	var session = r.mongoClient.NewSession()
	defer session.Close()

	_, err := session.
		DB(r.databaseName).
		C(consts.CollectionUsers).
		Upsert(bson.M{"UserId": user.UserId}, user)

	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(client mongo.Client, databaseName string) UserRepository {
	if err := client.EnsureIndex([]string{"UserId"}, true, "UserId", databaseName, consts.CollectionUsers); err != nil {
		panic(fmt.Sprintln("could not create a new client", err))
	}

	return &userRepository{mongoClient: client, databaseName: databaseName}
}
