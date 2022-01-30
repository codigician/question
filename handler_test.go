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
	"github.com/codigician/question/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFilterQuestions(t *testing.T) {
	mockService := mocks.NewMockService(gomock.NewController(t))
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
			scenario:           "Given valid query string it should call service with empty filter and return status ok",
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
	mockService := mocks.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario             string
		givenReqBody         interface{}
		expectedStatusCode   int
		mockErr              error
		expectedAlgoQuestion *q.Algorithm
	}{
		{
			scenario:           "Given bad request it should return 400",
			givenReqBody:       "invalid request body",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario:             "Given valid request body service call fails it should return 500",
			givenReqBody:         q.QuestionReqRes{},
			expectedAlgoQuestion: &q.Algorithm{},
			expectedStatusCode:   http.StatusInternalServerError,
			mockErr:              errors.New("an error"),
		},
		{
			scenario:             "Given valid request body service call succeeds it should return 201",
			givenReqBody:         q.QuestionReqRes{Title: "title", Content: "content"},
			expectedAlgoQuestion: &q.Algorithm{Title: "title", Content: "content"},
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
	mockService := mocks.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario           string
		givenQuestionID    string
		expectedStatusCode int
		mockErr            error
	}{
		{
			scenario:           "Given no question id it should return 404",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			scenario:           "Given valid question id it should return 200",
			givenQuestionID:    "1",
			expectedStatusCode: http.StatusOK,
		},
		{
			scenario:           "Given valid question id, service fails, it should return 500",
			givenQuestionID:    "2",
			expectedStatusCode: http.StatusInternalServerError,
			mockErr:            errors.New("an error"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			mockService.EXPECT().
				Get(gomock.Any(), tC.givenQuestionID).
				Return(&q.Algorithm{Tags: []string{"tag1", "tag2"}}, tC.mockErr)

			res, _ := http.Get(fmt.Sprintf("%s/questions/%s", srv.URL, tC.givenQuestionID))

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestUpdateQuestion(t *testing.T) {
	mockService := mocks.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario           string
		givenQuestionID    string
		givenQuestion      interface{}
		expectedStatusCode int
		expectedQuestion   *q.Algorithm
		mockErr            error
	}{
		{
			scenario:           "Given no question id it should return 404",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			scenario:           "Given question id but invalid request body it should return 400",
			givenQuestionID:    "1",
			givenQuestion:      "invalid request body",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			scenario:           "Given valid question id and valid request body it should return 200",
			givenQuestionID:    "2",
			givenQuestion:      q.QuestionReqRes{Title: "title", Content: "content"},
			expectedQuestion:   &q.Algorithm{Title: "title", Content: "content"},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			scenario:           "Given valid question id, valid request body, service fails it should return 500",
			givenQuestionID:    "3",
			givenQuestion:      q.QuestionReqRes{Title: "title", Content: "content"},
			expectedQuestion:   &q.Algorithm{Title: "title", Content: "content"},
			expectedStatusCode: http.StatusInternalServerError,
			mockErr:            errors.New("an error"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			mockService.EXPECT().
				Update(gomock.Any(), tC.givenQuestionID, tC.expectedQuestion).
				Return(tC.mockErr)

			bodyBytes, _ := json.Marshal(tC.givenQuestion)
			req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/questions/%s", srv.URL, tC.givenQuestionID), bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			res, _ := http.DefaultClient.Do(req)

			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestDeleteQuestion(t *testing.T) {
	mockService := mocks.NewMockService(gomock.NewController(t))
	srv := createTestServerAndRegisterRoutes(mockService)
	defer srv.Close()

	testCases := []struct {
		scenario           string
		givenQuestionID    string
		expectedStatusCode int
		mockErr            error
	}{
		{
			scenario:           "Given no question id it should return 404",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			scenario:           "Given valid question id delete successfully it should return 204",
			givenQuestionID:    "1",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			scenario:           "Given valid question id, service fails, it should return 500",
			givenQuestionID:    "2",
			expectedStatusCode: http.StatusInternalServerError,
			mockErr:            errors.New("an error"),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.scenario, func(t *testing.T) {
			mockService.EXPECT().
				Delete(gomock.Any(), tC.givenQuestionID).
				Return(tC.mockErr)

			req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/questions/%s", srv.URL, tC.givenQuestionID), nil)
			res, err := http.DefaultClient.Do(req)

			assert.Nil(t, err)
			assert.Equal(t, tC.expectedStatusCode, res.StatusCode)
		})
	}
}

func createTestServerAndRegisterRoutes(service *mocks.MockService) *httptest.Server {
	e := echo.New()
	handler := q.NewHandler(service)
	handler.RegisterRoutes(e)
	srv := httptest.NewServer(e)
	return srv
}
