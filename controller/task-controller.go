package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"graphql-pet/database"
	"graphql-pet/graph/model"
	"log"
	"time"
)

var mg *database.MongoInstance = database.GetMongoInstance()
var collection = mg.DB.Collection("tasks")

// GetAllTasks Функция получения списков всех таксов
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

// CreateTaskListing Функция для создания таска и добавления его в Базу Данных
func CreateTaskListing(taskInfo model.CreateTaskListingInput) *model.TaskListing {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inserted, err := collection.InsertOne(ctx, bson.M{
		"title":       taskInfo.Title,
		"description": taskInfo.Description,
		"company":     taskInfo.Company,
		"url":         taskInfo.URL,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedID := inserted.InsertedID.(primitive.ObjectID).Hex()
	taskListing := model.TaskListing{
		ID:          insertedID,
		Title:       taskInfo.Title,
		Description: taskInfo.Description,
		Company:     taskInfo.Company,
		URL:         taskInfo.URL,
	}
	return &taskListing
}

// GetTaskListing Получение таска
func GetTaskListing(id string) *model.TaskListing {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	filter := bson.M{"_id": taskId}

	var taskListing model.TaskListing
	err = collection.FindOne(ctx, filter).Decode(&taskListing)
	if err != nil {
		log.Fatal(err)
	}
	return &taskListing
}
