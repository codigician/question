package mongo

import (
	"context"

	"github.com/codigician/question"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_databaseListing    = "listing"
	_collectionQuestion = "question"
)

type (
	Mongo struct {
		uri    string
		client *mongo.Client
	}

	AlgoQuestion struct {
		ID         primitive.ObjectID `bson:"_id"`
		Title      string             `bson:"title"`
		Content    string             `bson:"content"`
		Template   string             `bson:"template"`
		Difficulty string             `bson:"difficulty"`
		Tags       []string           `bson:"tags"`
	}

	AlgoQuestions []AlgoQuestion
)

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

func (m *Mongo) Find(ctx context.Context, tags []string, difficulty question.Difficulty) ([]question.Algorithm, error) {
	filterQuery := bson.M{}
	if tags != nil {
		filterQuery["tags"] = bson.M{"$in": tags}
	}

	if difficulty != "" {
		filterQuery["difficulty"] = string(difficulty)
	}

	cursor, err := m.lq().Find(ctx, filterQuery)
	if err != nil {
		return nil, err
	}

	var dbQuestions AlgoQuestions
	err = cursor.All(ctx, &dbQuestions)
	return dbQuestions.to(), err
}

func (m *Mongo) Save(ctx context.Context, q *question.Algorithm) (string, error) {
	res, err := m.lq().InsertOne(ctx, fromQuestion(q))
	if err != nil {
		return "", err
	}
	id, _ := res.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (m *Mongo) lq() *mongo.Collection {
	return m.client.Database(_databaseListing).Collection(_collectionQuestion)
}

func (a *AlgoQuestion) to() question.Algorithm {
	return question.Algorithm{
		ID:         a.ID.Hex(),
		Title:      a.Title,
		Content:    a.Content,
		Template:   a.Template,
		Difficulty: question.Difficulty(a.Difficulty),
		Tags:       a.Tags,
	}
}

func (algoQuestions AlgoQuestions) to() (questions []question.Algorithm) {
	for _, aq := range algoQuestions {
		questions = append(questions, aq.to())
	}
	return questions
}

func fromQuestion(q *question.Algorithm) *AlgoQuestion {
	return &AlgoQuestion{
		ID:         primitive.NewObjectID(),
		Title:      q.Title,
		Content:    q.Content,
		Template:   q.Template,
		Difficulty: string(q.Difficulty),
		Tags:       q.Tags,
	}
}
