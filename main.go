package main

import (
	"context"
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
	// Altas cluster specifics. Be sure it includes
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
		fmt.Println("There was a problem connecting to your Atlas cluster. Check that the URI includes a validusername and password, and that your IP address has been whitelisted. Error: ")
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

	eloteIngredients := make([]string, 0)
	eloteIngredients = append(eloteIngredients, "corn")
	eloteIngredients = append(eloteIngredients, "mayonnaise")
	eloteIngredients = append(eloteIngredients, "cotija cheese")
	eloteIngredients = append(eloteIngredients, "sour cream")
	eloteIngredients = append(eloteIngredients, "lime")

	eloteRecipe := Recipe{"elote", eloteIngredients, 35}

	locoMocoIngredients := make([]string, 0)
	locoMocoIngredients = append(locoMocoIngredients, "ground beef")
	locoMocoIngredients = append(locoMocoIngredients, "butter")
	locoMocoIngredients = append(locoMocoIngredients, "onion")
	locoMocoIngredients = append(locoMocoIngredients, "egg")
	locoMocoIngredients = append(locoMocoIngredients, "bread bun")
	locoMocoIngredients = append(locoMocoIngredients, "mushrooms")

	locoMocoRecipe := Recipe{"loco moco", locoMocoIngredients, 54}

	patatasBravasIngredients := make([]string, 0)
	patatasBravasIngredients = append(patatasBravasIngredients, "potato")
	patatasBravasIngredients = append(patatasBravasIngredients, "tomato")
	patatasBravasIngredients = append(patatasBravasIngredients, "olive oil")
	patatasBravasIngredients = append(patatasBravasIngredients, "onion")
	patatasBravasIngredients = append(patatasBravasIngredients, "garlic")
	patatasBravasIngredients = append(patatasBravasIngredients, "paprika")

	patatasBravasRecipe := Recipe{"patas bravas", patatasBravasIngredients, 80}

	friedRiceIngredients := make([]string, 0)
	friedRiceIngredients = append(friedRiceIngredients, "rice")
	friedRiceIngredients = append(friedRiceIngredients, "soy sauce")
	friedRiceIngredients = append(friedRiceIngredients, "egg")
	friedRiceIngredients = append(friedRiceIngredients, "onion")
	friedRiceIngredients = append(friedRiceIngredients, "pea")
	friedRiceIngredients = append(friedRiceIngredients, "carrot")
	friedRiceIngredients = append(friedRiceIngredients, "sesame oil")

	friedRiceRecipe := Recipe{"fried rice", friedRiceIngredients, 40}

	// Create an interface of all the created recipes
	recipes := []interface{}{eloteRecipe, locoMocoRecipe, patatasBravasRecipe, friedRiceRecipe}
	insertManyResult, err := collection.InsertMany(context.TODO(), recipes)
	if err != nil {
		fmt.Println("Something went wrong trying to insert the new documents:")
		panic(err)
	}
	fmt.Println(insertManyResult.InsertedIDs, " documents successfully inserted.")

	/*
	 * *** FIND DOCUMENTS ***
	 *
	 * Now that we have data in Atlas, we can read it. To retrieve all of
	 * the data in a collection, we create a filter for recipes that take
	 * less than 45 minutes to prepare and sort by name (ascending)
	 */

	var filter = bson.M{"prepTimeInMinutes": bson.M{"$lt": 45}}
	options := options.Find()
	options.SetSort(bson.D{{"name", 1}}) // Sort by `name` field ascending

	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		fmt.Println("Something went wrong trying to find the documents:")
		panic(err)
	}

	// Loop through the returned recipes
	for cursor.Next(ctx) {
		// declare a result BSON object
		var result bson.M
		err := cursor.Decode(&result)
		// If there is a cursor.Decode error
		if err != nil {
			fmt.Println("cursor.Next() error:")
			panic(err)
		} else {
			var name = result["name"].(string)
			var ingredients = result["ingredients"]
			var prepTimeInMinutes = result["prepTimeInMinutes"].(int32)

			fmt.Println(name, "has the ingredients:", ingredients, "and takes", prepTimeInMinutes, "minutes to make")
		}
	}

	// We can also find a single document. Let's find the first document
	// that has the string "fried rice" in the name.
	var result Recipe
	var myFilter = bson.D{{"name", "fried rice"}}
	e := collection.FindOne(context.TODO(), myFilter).Decode(&result)
	if e != nil {
		fmt.Println("Something went wrong trying to find one document:")
		panic(e)
	}
	fmt.Println("found a document with the name fried rice", result)

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
	doc := bson.M{}
	decodeErr := myRes.Decode(&doc)
	if decodeErr != nil {
		fmt.Println("Something went wrong trying to decode the document:")
		panic(decodeErr)
	}
	fmt.Println("Here is the updated document:", doc)

	/*      *** DELETE DOCUMENTS ***
	 *
	 *      As with other CRUD methods, you can delete a single document
	 *      or all documents that match a specified filter. To delete all
	 *      of the documents in a collection, pass an empty filter to
	 *      the DeleteMany() method. In this example, we'll delete two of
	 *      the recipes.
	 */

	deletedRecipeNameList := make([]string, 0)
	deletedRecipeNameList = append(deletedRecipeNameList, "elotees")
	deletedRecipeNameList = append(deletedRecipeNameList, "fried rice")

	var deleteQuery = bson.M{"name": bson.M{"$in": deletedRecipeNameList}}
	deleteResult, err := collection.DeleteMany(context.TODO(), deleteQuery)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted %v documents in the recipes collection\n", deleteResult.DeletedCount)

}
