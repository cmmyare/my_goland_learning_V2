package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func InsertUser(user User) error {
	collection := MongoClient.Database(DB).Collection("users")
	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", inserted.InsertedID)
	return err
}

func FindUserByEmail(email string) (User, error) {
	var user User
	collection := MongoClient.Database(DB).Collection("users")
	filter := bson.M{"email": email}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
