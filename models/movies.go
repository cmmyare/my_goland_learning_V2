package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Movie  string             `json:"movie" bson:"movie"`
	Actors []string           `json:"actors" bson:"actors"`
}

func InserMovie(movie Movie) error {
	collection := MongoClient.Database(DB).Collection("movies")
	inserted, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", inserted.InsertedID)

	return err
}

func InsertMany(movies []Movie) error {
	newMovies := make([]interface{}, len(movies))
	for i, movie := range movies {
		newMovies[i] = movie
	}
	collection := MongoClient.Database(DB).Collection("movies")
	inserted, err := collection.InsertMany(context.TODO(), newMovies)
	if err != nil {
		return err
	}
	fmt.Println("Inserted multiple documents: ", inserted.InsertedIDs)
	return err
}
// func UpdateMovie(movieID string, movie Movie) error {
// 	id, err := primitive.ObjectIDFromHex(movieID)
// 	fmt.Println("Trying to update document with ID:", id.Hex())
// 	fmt.Println("Trying to update from req param:", movieID)
// 	if err != nil {
// 		return err
// 	}
// 	// filter := primitive.M{"_id": id}
// 	// update := primitive.M{"$set": movie}
// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": bson.M{
// 		"movie":  movie.Movie,
// 		"actors": movie.Actors,
// 	}}
// 	collection := MongoClient.Database(DB).Collection("movies")
// 	updated, err := collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	if updated.MatchedCount == 0 {
// 		// No document matched the filter (likely invalid/mismatched ObjectID)
// 		return fmt.Errorf("no document found with _id %s", id.Hex())
// 	}
// 	fmt.Println("Update attempt - matched:", updated.MatchedCount, "modified:", updated.ModifiedCount)
// 	return nil
// }

func UpdateMovie(movieID string, updateFields map[string]interface{}) error {
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}

	collection := MongoClient.Database(DB).Collection("movies")
	updated, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if updated.MatchedCount == 0 {
		return fmt.Errorf("no document found with _id %s", id.Hex())
	}
	return nil
}

func FindByName(movieName string) (Movie, error) {
	var result Movie
	filter := bson.M{"movie": movieName}
	collection := MongoClient.Database(DB).Collection("movies")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func FindByID(idStr string) (Movie, error) {
	var result Movie

	// Convert string to ObjectID
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return result, fmt.Errorf("invalid ObjectID: %v", err)
	}

	filter := bson.M{"_id": id}
	collection := MongoClient.Database(DB).Collection("movies")

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}


func FindAll(movieName string) []Movie {
	var results []Movie
	filter := bson.D{{Key: "movie", Value: movieName}}
	collection := MongoClient.Database(DB).Collection("movies")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	return results
}
func ListAll() ([]Movie, error) {
	var results []Movie
	collection := MongoClient.Database(DB).Collection("movies")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	return results, nil
}

func DeleteMovie(movieID string) error {
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	collection := MongoClient.Database(DB).Collection("movies")
	deleted, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	fmt.Println("Deleted a single document: ", deleted.DeletedCount)
	return nil
}

func DeleteAll() error {
	collection := MongoClient.Database(DB).Collection("movies")
	deleted, err := collection.DeleteMany(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		return err
	}
	fmt.Println("Deleted multiple documents: ", deleted.DeletedCount)
	return nil
}

func PartialUpdateMovie(movieID string, fields map[string]interface{}) error {
	id, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		return err
	}


	set := bson.M{}
	if v, ok := fields["movie"]; ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			set["movie"] = s
		}
	}
	if v, ok := fields["actors"]; ok {
		// Accept either []string or []interface{} convertible to []string
		switch arr := v.(type) {
		case []string:
			if len(arr) > 0 {
				set["actors"] = arr
			}
		case []interface{}:
			tmp := make([]string, 0, len(arr))
			for _, it := range arr {
				if s, ok := it.(string); ok {
					tmp = append(tmp, s)
				}
			}
			if len(tmp) > 0 {
				set["actors"] = tmp
			}
		}
	}

	if len(set) == 0 {
		// Nothing to update
		return fmt.Errorf("no valid fields provided to update")
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": set}
	collection := MongoClient.Database(DB).Collection("movies")
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no document found with _id %s", id.Hex())
	}
	return nil
}