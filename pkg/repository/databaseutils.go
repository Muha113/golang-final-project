package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Muha113/golang-final-project/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *UsersRepositoryInMemory) isExist(user model.User) bool {
	collection := u.db.Collection("Users")
	filter := bson.D{{Key: "useremail", Value: user.UserEmail}}
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
	filter = bson.D{{Key: "useremail", Value: user.UserEmail}}
	err = collection.FindOne(context.TODO(), filter).Decode(&res)
	if err == nil {
		return fmt.Errorf("Error %s", "duplicate email")
	}
	return nil
}
