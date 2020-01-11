package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepositoryInMemory struct {
	sync.RWMutex
	db *mongo.Database
}

func NewUsersRepositoryInMemory(database *mongo.Database) *UsersRepositoryInMemory {
	return &UsersRepositoryInMemory{
		db: database,
	}
}

func (u *UsersRepositoryInMemory) SaveUser(user model.User) error {
	u.RLock()
	err := u.checkModelForValid(user)
	if err != nil {
		return err
	}
	collection := u.db.Collection("Users")
	err = u.checkModelForExistence(user, collection)
	if err != nil {
		return err
	}
	u.RUnlock()
	u.Lock()
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	u.Unlock()
	return nil
}

func (u *UsersRepositoryInMemory) SaveTweet(tweet model.Tweet) error {
	return nil
}

func (u *UsersRepositoryInMemory) UpdateUser(user model.User) error {
	u.RLock()
	err := u.checkModelForValid(user)
	if err != nil {
		return err
	}
	collection := u.db.Collection("Users")
	err = u.checkModelForExistence(user, collection)
	if err != nil {
		return err
	}
	u.RUnlock()
	u.Lock()
	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": user.ID}, user)
	if err != nil {
		return err
	}
	u.Unlock()
	return nil
}

func (u *UsersRepositoryInMemory) getByID(id uint) (model.User, error) {

	return model.User{}, nil
}

func (u *UsersRepositoryInMemory) getByEmail(email string) (model.User, error) {
	return model.User{}, nil
}

func (u *UsersRepositoryInMemory) checkModelForValid(user model.User) error {
	val1, _ := regexp.MatchString("^\\w+@\\w+\\.[a-z]+", user.UserEmail)
	val2 := strings.ContainsRune(user.UserName, ' ')
	if !val1 || val2 || user.UserPasswordHash == "" {
		return fmt.Errorf("Error: %s\nJSON: %s", "bad input register json", user.ToString())
	}
	return nil
}

func (u *UsersRepositoryInMemory) checkModelForExistence(user model.User, coll *mongo.Collection) error {
	return nil
}
