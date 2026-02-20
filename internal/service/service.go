package service

import (
	"context"
	"net/http"
	"shortener/internal/logger"
	"shortener/internal/model"
	"shortener/internal/repository"
	"shortener/pkg/errormsg"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/helpers"
)

type Service struct {
	repo    repository.Repository
	baseURL string
}

func New(repo repository.Repository, baseURL string) *Service {
	return &Service{repo: repo, baseURL: baseURL}
}

type shortenRequest struct {
	URL         string `json:"url" binding:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty" binding:"omitempty,alphanum,max=20"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (s *Service) ShortenHandler(c *ginext.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errormsg.ErrorResponse{Error: true, Message: "invalid request: " + err.Error()})
		return
	}

	ctx := c.Request.Context()

	var shortURL string
	if req.CustomAlias != "" {
		exists, err := s.repo.LinkExists(ctx, req.CustomAlias)
		if err != nil {
			logger.Error("failed to check alias existence", "error", err)
			c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, errormsg.ErrorResponse{Error: true, Message: "alias already taken"})
			return
		}
		shortURL = req.CustomAlias
	} else {
		for {
			shortURL = helpers.CreateUUID()[:8]
			exists, err := s.repo.LinkExists(ctx, shortURL)
			if err != nil {
				logger.Error("failed to check short url existence", "error", err)
				c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
				return
			}
			if !exists {
				break
			}
		}
	}

	link := &model.Link{
		ShortURL:    shortURL,
		OriginalURL: req.URL,
		CustomAlias: req.CustomAlias,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.SaveLink(ctx, link); err != nil {
		logger.Error("failed to save link", "error", err)
		c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "failed to save link"})
		return
	}

	fullShortURL := s.baseURL + "/s/" + shortURL
	c.JSON(http.StatusCreated, shortenResponse{ShortURL: fullShortURL})
}

func (s *Service) RedirectHandler(c *ginext.Context) {
	shortURL := c.Param("short_url")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, errormsg.ErrorResponse{Error: true, Message: "short_url required"})
		return
	}

	ctx := c.Request.Context()

	link, err := s.repo.GetLink(ctx, shortURL)
	if err != nil {
		logger.Error("failed to get link", "error", err)
		c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
		return
	}
	if link == nil {
		c.JSON(http.StatusNotFound, errormsg.ErrorResponse{Error: true, Message: "link not found"})
		return
	}

	go func() {
		ctxBg := context.Background()
		analytics := &model.Analytics{
			ShortURL:  shortURL,
			Timestamp: time.Now(),
			UserAgent: c.GetHeader("User-Agent"),
			Referer:   c.GetHeader("Referer"),
		}
		if err := s.repo.SaveAnalytics(ctxBg, analytics); err != nil {
			logger.Error("failed to save analytics", "error", err, "short_url", shortURL)
		}
	}()

	c.Redirect(http.StatusFound, link.OriginalURL)
}

type analyticsItem struct {
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer"`
}

type analyticsResponse struct {
	ShortURL     string          `json:"short_url"`
	TotalClicks  int64           `json:"total_clicks"`
	RecentClicks []analyticsItem `json:"recent_clicks"`
}

func (s *Service) AnalyticsHandler(c *ginext.Context) {
	shortURL := c.Param("short_url")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, errormsg.ErrorResponse{Error: true, Message: "short_url required"})
		return
	}

	ctx := c.Request.Context()

	link, err := s.repo.GetLink(ctx, shortURL)
	if err != nil {
		logger.Error("failed to get link", "error", err)
		c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
		return
	}
	if link == nil {
		c.JSON(http.StatusNotFound, errormsg.ErrorResponse{Error: true, Message: "link not found"})
		return
	}

	total, err := s.repo.CountAnalytics(ctx, shortURL)
	if err != nil {
		logger.Error("failed to count analytics", "error", err)
		c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
		return
	}

	recent, err := s.repo.GetAnalytics(ctx, shortURL, 10) // последние 10
	if err != nil {
		logger.Error("failed to get recent analytics", "error", err)
		c.JSON(http.StatusInternalServerError, errormsg.ErrorResponse{Error: true, Message: "internal error"})
		return
	}

	items := make([]analyticsItem, len(recent))
	for i, a := range recent {
		items[i] = analyticsItem{
			Timestamp: a.Timestamp,
			UserAgent: a.UserAgent,
			Referer:   a.Referer,
		}
	}

	c.JSON(http.StatusOK, analyticsResponse{
		ShortURL:     shortURL,
		TotalClicks:  total,
		RecentClicks: items,
	})
}
