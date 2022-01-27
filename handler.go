package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	Service interface {
		Create(ctx context.Context, q *AlgorithmQuestion) (*AlgorithmQuestion, error)
		Filter(ctx context.Context, f Filter) ([]AlgorithmQuestion, error)
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

	q, err := h.qservice.Create(c.Request().Context(), req.To())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, CreateQuestionRes{q.ID})
}

func (h *Handler) FilterQuestions(c echo.Context) error {
	var filter Filter

	if tags := c.QueryParam("tags"); tags != "" {
		filter.Tags = strings.Split(tags, ",")
	}

	if difficulty := c.QueryParam("difficulty"); difficulty != "" {
		filter.Difficulty = Difficulty(difficulty)
	}

	questions, err := h.qservice.Filter(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, Questions(questions).To())
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

type Questions []AlgorithmQuestion

func (questions Questions) To() (filterRes []*QuestionReqRes) {
	for idx := range questions {
		filterRes = append(filterRes, FromQuestion(&questions[idx]))
	}
	return filterRes
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
