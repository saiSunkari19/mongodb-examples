package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URL = "mongodb://localhost:27017"

// const MONGO_URL = "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=myReplicaSet"

type User struct {
	Id   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
	Age  int64  `bson:"age" json:"age"`
}

func main() {
	db := GetDB()
	
	collection := db.Collection("users")
	
	// Setup change stream  options
	// Create pipeline
	pipeline := bson.A{
		bson.D{{"$project", bson.D{{"fullDocument", 1}}}},
	}
	changeStreamOptions := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	
	// Create a change stream
	chageStream, err := collection.Watch(context.Background(), pipeline, changeStreamOptions)
	if err != nil {
		log.Fatal(err)
	}
	
	// loop back through change stream and process the changes
	defer chageStream.Close(context.Background())
	for {
		hasNext := chageStream.Next(context.Background())
		if !hasNext {
			break
		}
		
		var chagedEvent bson.M
		if err := chageStream.Decode(&chagedEvent); err != nil {
			log.Panicln("Error decoding change event", err)
		} else {
			
			bz, err := json.Marshal(chagedEvent["fullDocument"])
			if err != nil {
				log.Panicln("Error marshal ", err)
			}
			var data User
			if err := json.Unmarshal(bz, &data); err != nil {
				log.Panicln("Error while Unmarshal ", err)
			}
			
			objectId := chagedEvent["fullDocument"].(bson.M)["_id"]
			// bz, err := json.Marshal(chagedEvent["fullDocument"])
			// if err != nil {
			// 	log.Panicln("Error marshal ", err)
			// }
			log.Println("ObjectId", objectId)
			log.Println("Data", data)
			log.Print("Change Event", fmt.Sprintf("%x", chagedEvent["fullDocument"].(bson.M)["_id"]))
		}
	}
	
	if err := chageStream.Err(); err != nil {
		log.Println("Change stream error", err)
	}
	
}

func GetDB() *mongo.Database {
	
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	
	return client.Database("mongo")
}
