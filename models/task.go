package models

import (
	"fmt"
	"time"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingInformation struct {
	ServiceType   string    `json:"service_type" bson:"service_type"`
	WorkType      string    `json:"work_type" bson:"work_type"`
	Instruction   string    `json:"instruction" bson:"instruction"`
	BookingDate   time.Time `json:"booking_date" bson:"booking_date"`
	ReachDate     string    `json:"reach_date" bson:"reach_date"`
	EstReachDate  time.Time `json:"est_reach_date" bson:"est_reach_date"`
	Location      string    `json:"location" bson:"location"`
	ServicePrice  float64   `json:"service_price" bson:"service_price"`
	UserLatLong   GeoPoint  `json:"user_latlong" bson:"user_latlong"`
}

type GeoPoint struct {
	Log float64 `json:"log" bson:"log"`
	Lat float64 `json:"lat" bson:"lat"`
}
type Task struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user" bson:"user"`
	TaskerID       primitive.ObjectID `json:"tasker" bson:"tasker"`
	CategoryID     primitive.ObjectID `json:"category" bson:"category"`
	ServiceLevelID primitive.ObjectID `json:"service_level_id" bson:"service_level_id"`
	BookingInfo    BookingInformation `json:"booking_information" bson:"booking_information"`
}

func InsertTask(task Task) error {
	collection := MongoClient.Database(DB).Collection("tasks")

	inserted, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		return err
	}

	fmt.Println("Inserted Task ID:", inserted.InsertedID)
	return nil
}


