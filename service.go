package question

import "context"

type (
	Repository interface {
		Get(ctx context.Context, id string) (*Algorithm, error)
		Save(ctx context.Context, q *Algorithm) (string, error)
		Find(ctx context.Context, tags []string, difficulty Difficulty) ([]Algorithm, error)
		Delete(ctx context.Context, id string) error
	}

	QuestionService struct {
		repository Repository
	}

	Filter struct {
		Tags       []string
		Difficulty Difficulty
	}
)

func NewService(repository Repository) *QuestionService {
	return &QuestionService{repository}
}

func (s *QuestionService) Create(ctx context.Context, q *Algorithm) (*Algorithm, error) {
	id, err := s.repository.Save(ctx, q)
	q.ID = id
	return q, err
}

func (s *QuestionService) Filter(ctx context.Context, f Filter) ([]Algorithm, error) {
	return s.repository.Find(ctx, f.Tags, f.Difficulty)
}

func (s *QuestionService) Get(ctx context.Context, id string) (*Algorithm, error) {
	return s.repository.Get(ctx, id)
}

func (s *QuestionService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}
