package question

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

type (
	Algorithm struct {
		ID         string
		Title      string
		Content    string
		Template   string
		Difficulty Difficulty

		Editorial Editorial

		Tags      []string
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
