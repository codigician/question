package question_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	q "github.com/codigician/question"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFilterQuestions(t *testing.T) {
	mockService := q.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario           string
		givenQueryString   string
		expectedFilter     q.Filter
		expectedStatusCode int
		mockErr            error
	}{
		{
			scenario:           "Given no query string it should call service with empty filter and return status ok",
			givenQueryString:   "?tags=tag1&difficulty=easy",
			expectedStatusCode: http.StatusOK,
			expectedFilter:     q.Filter{Tags: []string{"tag1"}, Difficulty: q.Easy},
		},
		{
			scenario:           "Given no query string it should call service with empty filter and return status ok",
			expectedStatusCode: http.StatusOK,
		},
		{
			scenario:           "Given valid query string it should call service when service fails it should return internal server error",
			givenQueryString:   "?difficulty=hard",
			mockErr:            errors.New("an error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedFilter:     q.Filter{Difficulty: q.Hard},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			mockService.EXPECT().
				Filter(gomock.Any(), tC.expectedFilter).
				Return(nil, tC.mockErr)

			res, err := http.Get(fmt.Sprintf("%s/questions%s", srv.URL, tC.givenQueryString))
			log.Println("error", err)

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestCreateQuestion(t *testing.T) {
	mockService := q.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario             string
		givenReqBody         interface{}
		expectedStatusCode   int
		mockErr              error
		expectedAlgoQuestion *q.AlgorithmQuestion
	}{
		{
			scenario:           "Given bad request it should return 400",
			givenReqBody:       "invalid request body",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario:             "Given valid request body service call fails it should return 500",
			givenReqBody:         q.QuestionReqRes{},
			expectedAlgoQuestion: &q.AlgorithmQuestion{},
			expectedStatusCode:   http.StatusInternalServerError,
			mockErr:              errors.New("an error"),
		},
		{
			scenario:             "Given valid request body service call succeeds it should return 201",
			givenReqBody:         q.QuestionReqRes{Title: "title", Content: "content"},
			expectedAlgoQuestion: &q.AlgorithmQuestion{Title: "title", Content: "content"},
			expectedStatusCode:   http.StatusCreated,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			mockService.EXPECT().
				Create(gomock.Any(), tC.expectedAlgoQuestion).
				Return(tC.expectedAlgoQuestion, tC.mockErr)

			bodyBytes, _ := json.Marshal(tC.givenReqBody)
			res, _ := http.Post(srv.URL+"/questions", "application/json", bytes.NewBuffer(bodyBytes))

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestGetQuestion(t *testing.T) {
}

func TestUpdateQuestion(t *testing.T) {
}

func TestDeleteQuestion(t *testing.T) {
}

func createTestServerAndRegisterRoutes(service *q.MockService) *httptest.Server {
	e := echo.New()
	handler := q.NewHandler(service)
	handler.RegisterRoutes(e)
	srv := httptest.NewServer(e)
	return srv
}
