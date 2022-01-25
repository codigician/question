package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_databaseListing    = "listing"
	_collectionQuestion = "question"
)

type Mongo struct {
	uri    string
	client *mongo.Client
}

func NewMongo(uri string) *Mongo {
	return &Mongo{uri: uri}
}

func (m *Mongo) Connect(ctx context.Context) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.uri))
	m.client = client
	return err
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *Mongo) Find(ctx context.Context, tags []string, difficulty Difficulty) (questions []AlgorithmQuestion, err error) {
	filterQuery := bson.M{
		"tags":       bson.M{"$in": tags},
		"difficulty": string(difficulty),
	}
	cursor, err := m.client.Database(_databaseListing).Collection(_collectionQuestion).Find(ctx, filterQuery)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &questions)
	return questions, err
}

func (m *Mongo) Save(ctx context.Context, q *AlgorithmQuestion) (string, error) {
	res, err := m.client.Database(_databaseListing).Collection(_collectionQuestion).InsertOne(ctx, q)
	if err != nil {
		return "", err
	}
	id, _ := res.InsertedID.(string)
	return id, nil
}
