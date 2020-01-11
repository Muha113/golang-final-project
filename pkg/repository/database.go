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
	db   *mongo.Database
	size uint
}

func NewUsersRepositoryInMemory(database *mongo.Database) *UsersRepositoryInMemory {
	ctn, _ := database.Collection("Users").CountDocuments(context.Background(), bson.D{})
	return &UsersRepositoryInMemory{
		db:   database,
		size: uint(ctn),
	}
}

func (u *UsersRepositoryInMemory) SaveUser(user model.User) error {
	u.RLock()
	err := u.checkModelForValid(user)
	if err != nil {
		u.RUnlock()
		return err
	}
	collection := u.db.Collection("Users")
	if u.size != 0 {
		err = u.checkUserModelForExistence(user)
		if err != nil {
			u.RUnlock()
			return err
		}
	}
	u.RUnlock()
	u.Lock()
	u.size++
	user.ID = u.size
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		u.size--
		u.Unlock()
		return err
	}
	u.Unlock()
	return nil
}

func (u *UsersRepositoryInMemory) UpdateUser(user model.User) error {
	u.RLock()
	err := u.checkModelForValid(user)
	if err != nil {
		return err
	}
	collection := u.db.Collection("Users")
	if u.isExist(user) {
		return fmt.Errorf("Error: %s", "user already exists")
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

// func (u *UsersRepositoryInMemory) getByID(id uint) (model.User, error) {
// 	collection := u.db.Collection("Users")
// 	filter := bson.D{{Key: "id", Value: id}}
// 	var res model.User
// 	err := collection.FindOne(context.TODO(), filter).Decode(&res)
// 	if err != nil {
// 		return model.User{}, err
// 	}
// 	return res, nil
// }

func (u *UsersRepositoryInMemory) GetUserByEmail(email string) (model.User, error) {
	collection := u.db.Collection("Users")
	filter := bson.D{{Key: "email", Value: email}}
	var res model.User
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (u *UsersRepositoryInMemory) isExist(user model.User) bool {
	collection := u.db.Collection("Users")
	filter := bson.D{{Key: "email", Value: user.UserEmail}}
	var res model.User
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return false
	}
	return true
}

func (u *UsersRepositoryInMemory) checkModelForValid(user model.User) error {
	val1, _ := regexp.MatchString("^\\w+@\\w+\\.[a-z]+", user.UserEmail)
	val2 := strings.ContainsRune(user.UserName, ' ')
	val3 := user.UserPasswordHash == ""
	if !val1 || val2 || val3 {
		return fmt.Errorf("Error: %s", "bad input register json")
	}
	return nil
}

func (u *UsersRepositoryInMemory) checkUserModelForExistence(user model.User) error {
	collection := u.db.Collection("Users")
	var res model.User
	filter := bson.D{{Key: "username", Value: user.UserName}}
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err == nil {
		return fmt.Errorf("Error: %s", "duplicate username")
	}
	filter = bson.D{{Key: "email", Value: user.UserEmail}}
	err = collection.FindOne(context.TODO(), filter).Decode(&res)
	if err == nil {
		return fmt.Errorf("Error %s", "duplicate email")
	}
	return nil
}
