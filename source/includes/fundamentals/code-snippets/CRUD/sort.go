package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	client.Database("tea").Collection("ratings").Drop(context.TODO())

	// begin insertDocs
	coll := client.Database("tea").Collection("ratings")
	docs := []interface{}{
		bson.D{{"type", "Masala"}, {"rating", 10}},
		bson.D{{"type", "Assam"}, {"rating", 5}},
		bson.D{{"type", "Oolong"}, {"rating", 7}},
		bson.D{{"type", "Earl Grey"}, {"rating", 8}},
		bson.D{{"type", "English Breakfast"}, {"rating", 5}},
	}

	result, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents inserted: %d\n", len(result.InsertedIDs))
	//end insertDocs

	fmt.Println("Ascending Sort:")
	{
		//begin ascending sort
		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"rating", 1}})

		cursor, err := coll.Find(context.TODO(), filter, opts)

		var results []bson.D
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
		//end ascending sort
	}

	fmt.Println("Descending Sort:")
	{
		//begin descending sort
		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"rating", -1}})

		cursor, err := coll.Find(context.TODO(), filter, opts)

		var results []bson.D
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
		//end descending sort
	}

	fmt.Println("Multi Sort:")
	{
		//begin multi sort
		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"rating", 1}, {"type", -1}})

		cursor, err := coll.Find(context.TODO(), filter, opts)

		var results []bson.D
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
		//end multi sort
	}

	fmt.Println("Aggregation Sort:")
	{
		// begin aggregate sort
		sortStage := bson.D{{"$sort", bson.D{{"rating", -1}, {"type", 1}}}}

		cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{sortStage})
		if err != nil {
			panic(err)
		}

		var results []bson.D
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
		// end aggregate sort
	}
}
