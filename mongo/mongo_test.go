package mongo_test

import (
	"context"
	"log"
	"testing"

	"github.com/codigician/question"
	qmongo "github.com/codigician/question/mongo"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	mongo     *qmongo.Mongo
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
	s.mongo = qmongo.NewMongo(uri)
	if err := s.mongo.Connect(ctx); err != nil {
		log.Fatalf("mongo connect: %v\n", err)
	}
}

func (s *QuestionMongoTestSuite) TearDownSuite() {
	ctx := context.Background()

	if err := s.mongo.Disconnect(ctx); err != nil {
		log.Fatalf("mongo disconnect: %v\n", err)
	}

	if err := s.container.Terminate(ctx); err != nil {
		log.Fatalf("container terminate: %v\n", err)
	}
}

func (s *QuestionMongoTestSuite) TestFind() {
	ctx := context.Background()

	s.insertQuestions(ctx,
		s.createMongoQuestion(question.Easy, []string{"binary tree", "tree", "data structures"}),
		s.createMongoQuestion(question.Easy, []string{"tree", "binary tree"}),
		s.createMongoQuestion(question.Medium, []string{"tree", "binary tree"}),
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

	actualQuestion := s.getQuestion(ctx, id)

	s.Nil(err)
	s.Equal(expectedQuestion.Content, actualQuestion.Content)
	s.Equal(expectedQuestion.Tags, actualQuestion.Tags)
	s.Equal(string(expectedQuestion.Difficulty), actualQuestion.Difficulty)
	s.Equal(expectedQuestion.Template, actualQuestion.Template)
	s.Equal(expectedQuestion.Title, actualQuestion.Title)
}

func (s *QuestionMongoTestSuite) TestGet() {
	ctx := context.Background()

	createdMongoQuestion := s.createMongoQuestion(question.Easy, []string{"binary tree", "tree", "data structures"})
	s.insertQuestions(ctx, createdMongoQuestion)

	question, err := s.mongo.Get(ctx, createdMongoQuestion.ID.Hex())
	s.Nil(err)
	s.Equal([]string{"binary tree", "tree", "data structures"}, question.Tags)
	s.Equal("easy", string(question.Difficulty))
}

func (s *QuestionMongoTestSuite) createMongoQuestion(diff question.Difficulty, tags []string) qmongo.AlgoQuestion {
	return qmongo.AlgoQuestion{
		ID:         primitive.NewObjectID(),
		Title:      "Title",
		Content:    "Content",
		Template:   "Template",
		Difficulty: string(diff),
		Tags:       tags,
	}
}

func (s *QuestionMongoTestSuite) createQuestion(diff question.Difficulty, tags []string) question.Algorithm {
	return question.Algorithm{
		Title:      "Title",
		Content:    "Content",
		Template:   "Template",
		Difficulty: diff,
		Tags:       tags,
	}
}

func (s *QuestionMongoTestSuite) insertQuestions(ctx context.Context, questions ...qmongo.AlgoQuestion) {
	var documents []interface{}
	for idx := range questions {
		documents = append(documents, questions[idx])
	}

	if _, err := s.client.Database(_database).Collection(_collection).InsertMany(ctx, documents); err != nil {
		log.Fatalf("insert many: %v\n", err)
	}
}

func (s *QuestionMongoTestSuite) getQuestion(ctx context.Context, id string) qmongo.AlgoQuestion {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("invalid object id")
	}

	res := s.client.Database(_database).Collection(_collection).FindOne(ctx, bson.M{"_id": oid})
	if res.Err() != nil {
		log.Fatal("find one: ", res.Err())
	}

	var actualQuestion qmongo.AlgoQuestion
	if err := res.Decode(&actualQuestion); err != nil {
		log.Fatal("actual question decode", err)
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
