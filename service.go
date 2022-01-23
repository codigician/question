package main

type (
	Repository interface {
		Save(q *AlgorithmQuestion) error
		Find(tags []string, difficulty Difficulty) ([]AlgorithmQuestion, error)
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

func (s *QuestionService) Create(q *AlgorithmQuestion) (*AlgorithmQuestion, error) {
	err := s.repository.Save(q)
	return q, err
}

func (s *QuestionService) Filter(f Filter) ([]AlgorithmQuestion, error) {
	return s.repository.Find(f.Tags, f.Difficulty)
}
