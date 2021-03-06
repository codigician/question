package question

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	Service interface {
		Get(ctx context.Context, id string) (*Algorithm, error)
		Create(ctx context.Context, q *Algorithm) (*Algorithm, error)
		Filter(ctx context.Context, f Filter) ([]Algorithm, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, q *Algorithm) error
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
		log.Printf("create question: %v\n", err)
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
		log.Printf("filter questions: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, Questions(questions).To())
}

func (h *Handler) GetQuestion(c echo.Context) error {
	id := c.Param("id")

	q, err := h.qservice.Get(c.Request().Context(), id)
	if err != nil {
		log.Printf("get question: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, FromQuestion(q))
}

func (h *Handler) UpdateQuestion(c echo.Context) error {
	id := c.Param("id")

	var req QuestionReqRes
	if err := c.Bind(&req); err != nil {
		return err
	}

	err := h.qservice.Update(c.Request().Context(), id, req.To())
	if err != nil {
		log.Printf("update question: %v\n", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteQuestion(c echo.Context) error {
	id := c.Param("id")

	if err := h.qservice.Delete(c.Request().Context(), id); err != nil {
		log.Printf("delete question: %v\n", err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type Questions []Algorithm

func (questions Questions) To() (filterRes []*QuestionReqRes) {
	for idx := range questions {
		filterRes = append(filterRes, FromQuestion(&questions[idx]))
	}
	return filterRes
}

func FromQuestion(q *Algorithm) *QuestionReqRes {
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

func (r QuestionReqRes) To() *Algorithm {
	return &Algorithm{
		Title:      r.Title,
		Content:    r.Content,
		Template:   r.Template,
		Difficulty: Difficulty(r.Difficulty),
		Tags:       r.Tags,
	}
}
