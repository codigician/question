package question_test

import (
	"context"
	"testing"

	"github.com/codigician/question"
	"github.com/codigician/question/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate_RepositoryReturnsID_GetSameQuestionID(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return("1", nil)

	service := question.NewService(mockRepository)

	q, err := service.Create(context.Background(), &question.Algorithm{})

	assert.Nil(t, err)
	assert.Equal(t, "1", q.ID)
}

func TestCreate_RepositoryReturnsErr_ReturnErr(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return("", assert.AnError)

	service := question.NewService(mockRepository)

	_, err := service.Create(context.Background(), &question.Algorithm{})

	assert.NotNil(t, err)
}

func TestFilter_GivenFilter_ExpectRepositoryCallWithFilters(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().Find(gomock.Any(), []string{"tree"}, question.Difficulty("hard")).
		Return([]question.Algorithm{}, nil)

	service := question.NewService(mockRepository)

	res, err := service.Filter(context.Background(), question.Filter{
		Tags:       []string{"tree"},
		Difficulty: "hard",
	})

	assert.Nil(t, err)
	assert.Equal(t, 0, len(res))
}

func TestGet_GivenID_CallRepository(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().Get(gomock.Any(), "1").
		Return(&question.Algorithm{}, nil)

	service := question.NewService(mockRepository)

	_, err := service.Get(context.Background(), "1")

	assert.Nil(t, err)
}

func TestDelete_GivenID_CallRepository(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockRepository.EXPECT().Delete(gomock.Any(), "1").Return(nil)

	service := question.NewService(mockRepository)

	err := service.Delete(context.Background(), "1")

	assert.Nil(t, err)
}

func TestUpdate_GivenIDAndQuestion_CallRepository(t *testing.T) {
	mockRepository := mocks.NewMockRepository(gomock.NewController(t))
	mockQuestion := &question.Algorithm{
		Title:      "Title",
		Content:    "Content",
		Template:   "Templaet",
		Difficulty: "hard",
	}
	mockRepository.EXPECT().Update(gomock.Any(), "1", mockQuestion).
		Return(nil)

	service := question.NewService(mockRepository)

	err := service.Update(context.Background(), "1", mockQuestion)

	assert.Nil(t, err)
}
