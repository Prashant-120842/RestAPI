package Client

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Property struct {
	Propid       int
	Propertyname string
	Address      string
	City         string `json:"City"`
	Bedrooms     int    `json:"Bedrooms"`
}

func GetPropertiesById(propid int) Property {
	var prp Property
	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties") //from given db access books collection(Table)

	filter := bson.M{"propid": propid}

	err := collection.FindOne(context.TODO(), filter).Decode(&prp)

	if err != nil {
		log.Println(err)
	}
	return prp
}

func FilterProperties(param1 string, param2 int) []Property {
	var Properties []Property

	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties")

	filter := bson.M{"city": param1, "bedrooms": param2}

	cur, _ := collection.Find(context.TODO(), filter)

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var prop Property

		err := cur.Decode(&prop)

		if err != nil {
			log.Fatal(err)
		}

		Properties = append(Properties, prop)

	}
	return Properties

}

func InsertProperty(prop Property) {
	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties")

	result, err := collection.InsertOne(context.TODO(), prop)

	log.Println(result)
	log.Println(err)
}

func GetProperties() []Property {
	var Properties []Property

	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties")
	fmt.Println(collection)

	cur, _ := collection.Find(context.TODO(), bson.M{})

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var prop Property

		err := cur.Decode(&prop)

		if err != nil {
			log.Fatal(err)
		}

		Properties = append(Properties, prop)

	}
	//fmt.Println("Properties", Properties)
	return Properties

}

func PutProperty(propid int, prop Property) {

	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties")

	filter := bson.M{"propid": propid}

	update := bson.M{"propid": prop.Propid, "propertyname": prop.Propertyname, "address": prop.Address, "city": prop.City, "bedrooms": prop.Bedrooms}

	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": update})

	if err != nil {
		log.Fatal(err.Error())
	}

	var updatedProperty Property

	if result.MatchedCount == 1 {
		err := collection.FindOne(context.TODO(), bson.M{"propid": prop.Propid}).Decode(&updatedProperty)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	log.Println(updatedProperty)

}

func DeleteBooksById(propid int) {

	session := Connect()

	collection := session.Database("RestfulAPI").Collection("properties")

	filter := bson.M{"propid": propid}

	result, err := collection.DeleteOne(context.TODO(), filter)

	fmt.Println(result, err)

}

func Connect() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB Connected")
	return client
}
