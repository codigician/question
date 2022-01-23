package main

type (
	Difficulty string
	Tag        string
)

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type (
	AlgorithmQuestion struct {
		ID         string
		Title      string
		Content    string
		Template   string
		Difficulty Difficulty

		Editorial Editorial

		Tags      []Tag
		TestCases []TestCase
	}

	TestCase struct {
		Input  string
		Output string
	}

	Editorial struct {
		Explanation string
	}
)
