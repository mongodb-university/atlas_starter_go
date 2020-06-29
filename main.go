package main

import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "context"
    "time"
    "log"
    "fmt"
)

type Task struct {
    Name string
    Description  string
}

func main(){


    // CONNECT TO YOUR ATLAS CLUSTER:
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(
      "mongodb+srv://<username>:<password>@<cluster-name>/test?retryWrites=true&w=majority"
    ))
    if err != nil { log.Fatal(err) }

    err = client.Ping(ctx, nil)
    
    if err != nil {
        log.Fatal("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been whitelisted. Error: ", err)
    }

    fmt.Println("Connected to MongoDB!")

    // INSERT TASK DOCUMENTS: 

    collection := client.Database("tutorials").Collection("tasks")

    groceriesTask := Task { "Buy Groceries", "Milk, Eggs, Bread"}
    workoutTask := Task { "Workout", "Pullups, Pushups"}
    homeworkTask := Task { "Homework", "Read chapter 13 - 15"}
    meditationTask := Task { "Meditation", "Meditate 20 minutes"}

    tasks := []interface{}{ groceriesTask, workoutTask, homeworkTask, meditationTask }
    insertManyResult, err := collection.InsertMany(context.TODO(), tasks)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

    // READ DOCUMENTS: 

    cursor, err := collection.Find(context.TODO(), bson.D{})    
    for cursor.Next(ctx) {
        // declare a result BSON object
        var result bson.M
        err := cursor.Decode(&result)
        // If there is a cursor.Decode error
        if err != nil {
            log.Fatal("cursor.Next() error:", err)
        } else {
            fmt.Println("result:", result)
        }
    }

    // UPDATE DOCUMENT: 
    filter := bson.D{{"name", "Buy Groceries"}}
    update := bson.D{
        {"$set", bson.D{
            {"description", "Milk, Eggs, Bread, Cheese"},
        }},
    }

    updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)


    // READ ONE DOCUMENT: 
    // create a value into which the result can be decoded
    var myTask Task

    err = collection.FindOne(context.TODO(), filter).Decode(&myTask)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found a single document: %+v\n", myTask)


    // DELETE ALL DOCUMENTS: 
    deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Deleted %v documents in the tasks collection\n", deleteResult.DeletedCount)

}