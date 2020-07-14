package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A Recipe Struct allows you to insert recipe documents into your
// collection

type Recipe struct {
	Name              string
	Ingredients       []string
	PrepTimeInMinutes int `json:"prepTimeInMinutes" bson:"prepTimeInMinutes"`
}

func main() {

	// TODO:
	// Replace the placeholder connection string below with your
	// Atlas cluster specifics. Be sure it includes
	// a valid username and password! Note that in a production environment,
	// you do not want to store your password in plain-text here.
	var mongoUri = "<Your Atlas Connection String>"

	// CONNECT TO YOUR ATLAS CLUSTER:
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoUri,
	))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, nil)

	if err != nil {
		fmt.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a valid username and password, and that your IP address has been added to the access list. Error: ")
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Provide the name of the database and collection you want to use.
	// If they don't already exist, the driver and Atlas will create them
	// automatically when you first write data.
	var dbName = "myDatabase"
	var collectionName = "recipes"
	collection := client.Database(dbName).Collection(collectionName)

	/*      *** INSERT DOCUMENTS ***
	 *
	 * You can insert individual documents using collection.Insert().
	 * In this example, we're going to create 4 documents and then
	 * insert them all in one call with InsertMany().
	 */

	eloteRecipe := Recipe{
		Name:              "elote",
		Ingredients:       []string{"corn", "mayonnaise", "cotija cheese", "sour cream", "lime"},
		PrepTimeInMinutes: 35,
	}

	locoMocoRecipe := Recipe{
		Name:              "loco moco",
		Ingredients:       []string{"ground beef", "butter", "onion", "egg", "bread bun", "mushrooms"},
		PrepTimeInMinutes: 54,
	}

	patatasBravasRecipe := Recipe{
		Name:              "patatas bravas",
		Ingredients:       []string{"potato", "tomato", "olive oil", "onion", "garlic", "paprika"},
		PrepTimeInMinutes: 80,
	}

	friedRiceRecipe := Recipe{
		Name:              "fried rice",
		Ingredients:       []string{"rice", "soy sauce", "egg", "onion", "pea", "carrot", "sesame oil"},
		PrepTimeInMinutes: 40,
	}

	// Create an interface of all the created recipes
	recipes := []interface{}{eloteRecipe, locoMocoRecipe, patatasBravasRecipe, friedRiceRecipe}
	insertManyResult, err := collection.InsertMany(context.TODO(), recipes)
	if err != nil {
		fmt.Println("Something went wrong trying to insert the new documents:")
		panic(err)
	}

	fmt.Println(len(insertManyResult.InsertedIDs), "documents successfully inserted.")

	/*
	 * *** FIND DOCUMENTS ***
	 *
	 * Now that we have data in Atlas, we can read it. To retrieve all of
	 * the data in a collection, we create a filter for recipes that take
	 * less than 45 minutes to prepare and sort by name (ascending)
	 */

	var filter = bson.M{"prepTimeInMinutes": bson.M{"$lt": 45}}
	options := options.Find()

	// Sort by `name` field ascending
	options.SetSort(bson.D{{"name", 1}})

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		fmt.Println("Something went wrong trying to find the documents:")
		panic(err)
	}

	defer func() {
		cursor.Close(context.Background())
	}()

	// Loop through the returned recipes
	for cursor.Next(ctx) {
		recipe := Recipe{}
		err := cursor.Decode(&recipe)

		// If there is an error decoding the cursor into a Recipe
		if err != nil {
			fmt.Println("cursor.Next() error:")
			panic(err)
		} else {
			fmt.Println(recipe.Name, "has", len(recipe.Ingredients), "ingredients, and takes", recipe.PrepTimeInMinutes, "minutes to make.")
		}
	}

	// We can also find a single document. Let's find the first document
	// that has the string "fried rice" in the name.
	var result Recipe
	var myFilter = bson.D{{"ingredients", "potato"}}
	e := collection.FindOne(context.TODO(), myFilter).Decode(&result)
	if e != nil {
		fmt.Println("Something went wrong trying to find one document:")
		panic(e)
	}
	fmt.Println("Found a document with the ingredient potato", result)

	/*
	 * *** UPDATE A DOCUMENT ***
	 *
	 * You can update a single document or multiple documents in a single call.
	 *
	 * Here we update the prepTimeInMinutes value on the document we
	 * just found.
	 */

	var updateDoc = bson.D{{"$set", bson.D{{"prepTimeInMinutes", 72}}}}
	var myRes = collection.FindOneAndUpdate(ctx, myFilter, updateDoc, nil)
	if myRes.Err() != nil {
		fmt.Println("Something went wrong trying to update one document:")
		panic(myRes.Err())
	}

	_recipe := Recipe{}
	decodeErr := myRes.Decode(&_recipe)
	if decodeErr != nil {
		fmt.Println("Something went wrong trying to decode the document:")
		panic(decodeErr)
	}

	// indent the Recipe output to cleanly print the document using json.MarshallIndent
	updatedRecipe, _ := json.MarshalIndent(_recipe, "", "\t")
	fmt.Println("The following document has been updated: \n", string(updatedRecipe))

	/*      *** DELETE DOCUMENTS ***
	 *
	 *      As with other CRUD methods, you can delete a single document
	 *      or all documents that match a specified filter. To delete all
	 *      of the documents in a collection, pass an empty filter to
	 *      the DeleteMany() method. In this example, we'll delete two of
	 *      the recipes.
	 */

	deletedRecipeNameList := [...]string{"elote", "fried rice"}

	var deleteQuery = bson.M{"name": bson.M{"$in": deletedRecipeNameList}}
	deleteResult, err := collection.DeleteMany(context.TODO(), deleteQuery)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted", deleteResult.DeletedCount, "documents in the recipes collection\n")

}
