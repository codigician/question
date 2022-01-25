package main_test

import (
	"context"
	"log"
	"testing"

	question "github.com/codigician/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_database   = "listing"
	_collection = "question"
)

type QuestionMongoTestSuite struct {
	suite.Suite
	container testcontainers.Container
	client    *mongo.Client
	mongo     *question.Mongo
}

func TestIntegrationQuestionMongo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping mongodb integration tests")
	}

	suite.Run(t, new(QuestionMongoTestSuite))
}

func (s *QuestionMongoTestSuite) SetupSuite() {
	var (
		ctx = context.Background()
		uri = "mongodb://localhost:27017"
	)

	s.container = s.createMongoDBContainer(ctx)
	s.client = s.createMongoDBClient(ctx, uri)
	s.mongo = question.NewMongo(uri)
	s.mongo.Connect(ctx)
}

func (s *QuestionMongoTestSuite) TearDownSuite() {
	ctx := context.Background()

	s.mongo.Disconnect(ctx)
	s.container.Terminate(ctx)
}

func (s *QuestionMongoTestSuite) TestFind() {
	ctx := context.Background()

	s.insertQuestions(ctx,
		s.createQuestion(question.Easy, []string{"binary tree", "tree", "data structures"}),
		s.createQuestion(question.Easy, []string{"tree", "binary tree"}),
		s.createQuestion(question.Medium, []string{"tree", "binary tree"}),
	)

	questions, err := s.mongo.Find(ctx, []string{"tree", "binary tree"}, question.Easy)
	log.Println(questions)

	s.Nil(err)
	s.Len(questions, 2)
}

func (s *QuestionMongoTestSuite) TestSave() {
	ctx := context.Background()

	expectedQuestion := s.createQuestion(question.Hard, []string{"data structures"})

	id, err := s.mongo.Save(ctx, &expectedQuestion)

	s.Nil(err)
	s.Equal(expectedQuestion, s.getQuestion(ctx, id))
}

func (s *QuestionMongoTestSuite) createQuestion(diff question.Difficulty, tags []string) question.AlgorithmQuestion {
	return question.AlgorithmQuestion{
		// TODO: make this id created by mongodb
		ID:         uuid.NewString(),
		Title:      "Title",
		Content:    "Content",
		Template:   "Template",
		Difficulty: diff,
		Tags:       tags,
	}
}

func (s *QuestionMongoTestSuite) insertQuestions(ctx context.Context, questions ...question.AlgorithmQuestion) {
	var documents []interface{}
	for idx := range questions {
		documents = append(documents, questions[idx])
	}

	s.client.Database(_database).Collection(_collection).InsertMany(ctx, documents)
}

func (s *QuestionMongoTestSuite) getQuestion(ctx context.Context, id string) question.AlgorithmQuestion {
	res := s.client.Database(_database).Collection(_collection).FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		log.Fatal(res.Err())
	}

	var actualQuestion question.AlgorithmQuestion
	if err := res.Decode(&actualQuestion); err != nil {
		log.Fatal(err)
	}

	return actualQuestion
}

func (s *QuestionMongoTestSuite) createMongoDBClient(ctx context.Context, mongoURI string) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func (s *QuestionMongoTestSuite) createMongoDBContainer(ctx context.Context) testcontainers.Container {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{"27017:27017"},
			// WaitingFor:   wait.ForLog(""),
		},
		Started: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return container
}
