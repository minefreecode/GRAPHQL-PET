package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"graphql-pet/database"
	"graphql-pet/graph/model"
	"log"
	"time"
)

var mg *database.MongoInstance = database.GetMongoInstance()
var collection = mg.DB.Collection("tasks")

func GetAllTasks() []*model.TaskListing {
	query := bson.D{{}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	tasks := make([]*model.TaskListing, 0)
	err = cursor.All(context.TODO(), &tasks)
	if err != nil {
		log.Fatal(err)
	}
	return tasks
}
