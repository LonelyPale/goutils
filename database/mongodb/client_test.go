// Created by LonelyPale at 2019-12-07

package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"testing"
)

const configPath = "./config/mongodb.conf.toml"

var client *Client

func TestInsertOne(t *testing.T) {
	collection := client.Database("TestDB").Collection("test")

	id, err := collection.InsertOne(nil, map[string]string{"title": "hello world!"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(id)
}

func TestFindOne(t *testing.T) {
	collection := client.Database("TestDB").Collection("test")

	filter := map[string]string{}
	result := bson.M{}
	err := collection.FindOne(nil, &result, filter)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestMain(m *testing.M) {
	clientOptions := NewClientOptionsFromConfig(configPath)

	var err error
	client, err = Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()

	err = client.Disconnect()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}
