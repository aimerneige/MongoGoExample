package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {

	// Database config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://<username>:<password>@accounts.cojaq.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	// Close the connection once finished
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connect to mongodb successful!")

	// Print all databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All databases: ", databases)

	collection := client.Database("test").Collection("students")

	// Some data
	ZhangSan := Student{
		Name: "ZhangSan",
		Age:  18,
		City: "Shanghai",
	}
	LiShi := Student{
		Name: "LiShi",
		Age:  19,
		City: "Beijing",
	}
	WangWu := Student{
		Name: "WangWu",
		Age:  20,
		City: "Beijing",
	}

	// Insert single document
	insertResult, err := collection.InsertOne(context.TODO(), ZhangSan)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println(insertResult)

	// Insert multiple documents
	students := []interface{}{LiShi, WangWu}

	insertManyResult, err := collection.InsertMany(context.TODO(), students)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple document: ", insertManyResult.InsertedIDs)
	fmt.Println(insertManyResult)

	// Update a document
	// filter := bson.D{{"name", "ZhangSan"}}
	filter := bson.D{
		{
			"name", "ZhangSan",
		},
	}

	// update := bson.D{{"$inc", bson.D{{"age", 40}}}}
	update := bson.D{
		{
			"$inc", bson.D{
				{
					"age", 1,
				},
			},
		},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find a single document
	var result Student

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)

	// Find multiple documents returns a cursor
	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*Student

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	// Delete all the documents in the collection
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}
