package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//UsersRepositoryInMemory : represents db functionality
type UsersRepositoryInMemory struct {
	sync.RWMutex
	db   *mongo.Database
	size uint
}

//NewUsersRepositoryInMemory : creates UsersRepositoryInMemory object
func NewUsersRepositoryInMemory(database *mongo.Database) *UsersRepositoryInMemory {
	ctn, _ := database.Collection("Users").CountDocuments(context.Background(), bson.D{})
	return &UsersRepositoryInMemory{
		db:   database,
		size: uint(ctn),
	}
}

//GetSize : return documents quontity in collection
func (u *UsersRepositoryInMemory) GetSize() uint {
	return u.size
}

//SaveUser : save user in db if there is no such user, else return error
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

//UpdateUser : update user in db if this user exists, other case return error
func (u *UsersRepositoryInMemory) UpdateUser(user model.User) error {
	u.RLock()
	err := u.checkModelForValid(user)
	if err != nil {
		return err
	}
	collection := u.db.Collection("Users")
	if !u.isExist(user) {
		return fmt.Errorf("Error: %s", "user already exists")
	}
	u.RUnlock()
	u.Lock()
	filter := bson.D{{Key: "useremail", Value: user.UserEmail}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "userfollowing", Value: user.UserFollowing}, {Key: "usertweets", Value: user.UserTweets}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	u.Unlock()
	return nil
}

//GetUserByUserName : return user was found by name, if there is no such user return error
func (u *UsersRepositoryInMemory) GetUserByUserName(userName string) (model.User, error) {
	u.RLock()
	defer u.RUnlock()
	collection := u.db.Collection("Users")
	filter := bson.D{{Key: "username", Value: userName}}
	var res model.User
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

//GetUserByEmail : return user was found by email, if there is no such user return error
func (u *UsersRepositoryInMemory) GetUserByEmail(email string) (model.User, error) {
	u.RLock()
	defer u.RUnlock()
	collection := u.db.Collection("Users")
	filter := bson.D{{Key: "useremail", Value: email}}
	var res model.User
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}
