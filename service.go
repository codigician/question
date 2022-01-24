package main

import "context"

type (
	Repository interface {
		Save(ctx context.Context, q *AlgorithmQuestion) (string, error)
		Find(ctx context.Context, tags []string, difficulty Difficulty) ([]AlgorithmQuestion, error)
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

func (s *QuestionService) Create(ctx context.Context, q *AlgorithmQuestion) (*AlgorithmQuestion, error) {
	id, err := s.repository.Save(ctx, q)
	q.ID = id
	return q, err
}

func (s *QuestionService) Filter(ctx context.Context, f Filter) ([]AlgorithmQuestion, error) {
	return s.repository.Find(ctx, f.Tags, f.Difficulty)
}
