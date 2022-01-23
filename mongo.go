package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *Mongo) Collection(db, collection string) *mongo.Collection {
	return m.client.Database(db).Collection(collection)
}

func (m *Mongo) Find(tags []string, difficulty Difficulty) ([]AlgorithmQuestion, error) {
	return nil, nil
}

func (m *Mongo) Save(q *AlgorithmQuestion) error {
	return nil
}
