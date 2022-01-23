package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	Service interface {
		Create(q *AlgorithmQuestion) (*AlgorithmQuestion, error)
		Filter(f Filter) ([]AlgorithmQuestion, error)
	}

	Handler struct {
		qservice Service
	}

	CreateQuestionRes struct {
		ID string `json:"id"`
	}

	QuestionReqRes struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Template   string   `json:"template"`
		Difficulty string   `json:"difficulty"`
		Tags       []string `json:"tags"`
	}
)

func NewHandler(qservice Service) *Handler {
	return &Handler{qservice}
}

func (h *Handler) RegisterRoutes(router *echo.Echo) {
	// filtering:  /questions?tag=trees&tag=bfs&tag=dfs&difficulty=easy
	router.GET("/questions", h.FilterQuestions)

	router.GET("/questions/:id", h.GetQuestion)

	router.POST("/questions", h.CreateQuestion)

	router.PUT("/questions/:id", h.UpdateQuestion)
	router.DELETE("/questions/:id", h.DeleteQuestion)
}

func (h *Handler) CreateQuestion(c echo.Context) error {
	var req QuestionReqRes
	if err := c.Bind(&req); err != nil {
		return err
	}

	q, err := h.qservice.Create(req.To())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, CreateQuestionRes{q.ID})
}

func (h *Handler) FilterQuestions(c echo.Context) error {
	tags := c.Param("tags")
	difficulty := c.Param("difficulty")

	f := Filter{
		Tags:       strings.Split(tags, ","),
		Difficulty: Difficulty(difficulty),
	}
	questions, err := h.qservice.Filter(f)
	if err != nil {
		return err
	}

	var filterRes []*QuestionReqRes
	for idx := range questions {
		filterRes = append(filterRes, FromQuestion(&questions[idx]))
	}
	return c.JSON(http.StatusOK, filterRes)
}

func (h *Handler) GetQuestion(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	return c.String(http.StatusNotImplemented, "not implemented")
}

func (h *Handler) UpdateQuestion(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	return c.String(http.StatusNotImplemented, "not implemented")
}

func (h *Handler) DeleteQuestion(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	return c.String(http.StatusNotImplemented, "not implemented")
}

func FromQuestion(q *AlgorithmQuestion) *QuestionReqRes {
	var tags []string
	for _, tag := range q.Tags {
		tags = append(tags, string(tag))
	}

	return &QuestionReqRes{
		Title:      q.Title,
		Content:    q.Content,
		Template:   q.Template,
		Difficulty: string(q.Difficulty),
		Tags:       tags,
	}
}

func (r QuestionReqRes) To() *AlgorithmQuestion {
	return &AlgorithmQuestion{
		Title:      r.Title,
		Content:    r.Content,
		Template:   r.Template,
		Difficulty: Difficulty(r.Difficulty),
		Tags:       r.Tags,
	}
}
