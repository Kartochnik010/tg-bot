package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kartochnik010/tg-bot/internal/pkg/logger"
)

func NewCatsApiService(client *http.Client) *CatsApiService {
	return &CatsApiService{
		c: client,
	}
}

type CatsApiService struct {
	c     *http.Client
	cache []CatApiResponse
}

type CatPicture struct {
	Breeds []any  `json:"breeds"`
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CatApiResponse struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Temperament  string `json:"temperament"`
	ID           string `json:"id"`
	WikipediaURL string `json:"wikipedia_url,omitempty"`
	Weight       struct {
		Metric string `json:"metric"`
	} `json:"weight"`
	Origin   string `json:"origin"`
	LifeSpan string `json:"life_span"`
}

func (c CatApiResponse) String() string {
	return fmt.Sprintf(`%s

%v

Temperament is %v
Originated from %v
Expected life span is %v years
Learn more about %v at %v
`, c.Name, c.Description, c.Temperament, c.Origin, c.LifeSpan, c.Name, c.WikipediaURL)
}
func (s *CatsApiService) GetRandomCat(ctx context.Context) (*CatPicture, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.GetRandomCatInfo")

	var cats []CatPicture
	resp, err := s.c.Get("https://api.thecatapi.com/v1/images/search")
	if err != nil {
		// log.WithError(err).Error("failed to fetch music")
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.WithError(err).Error("failed to read body")
		return nil, err
	}
	if err = json.Unmarshal(body, &cats); err != nil {
		// log.WithError(err).Error("failed to unmarshal")
		return nil, err
	}

	if len(cats) == 0 {
		// log.Error("no cats found")
		return nil, fmt.Errorf("no cats found")
	}

	return &cats[0], nil
}

func (s *CatsApiService) GetAllBreeds(ctx context.Context) ([]CatApiResponse, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.GetAllBreeds")
	if len(s.cache) != 0 {
		return s.cache, nil
	}
	var cats []CatApiResponse
	resp, err := s.c.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		// log.WithError(err).Error("failed to fetch breeds")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.WithError(err).Error("failed to read body")
		return nil, err
	}
	if err = json.Unmarshal(body, &cats); err != nil {
		// log.WithError(err).Error("failed to unmarshal")
		return nil, err
	}

	if len(cats) == 0 {
		// log.Error("no cats found")
		return nil, fmt.Errorf("no cats found")
	}

	s.cache = cats

	return cats, nil
}

func (s *CatsApiService) GetBreed(ctx context.Context, breed string) (*CatApiResponse, error) {
	log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.GetAllBreeds")

	var cat CatApiResponse
	resp, err := s.c.Get("https://api.thecatapi.com/v1/breeds/" + breed)
	if err != nil {
		// log.WithError(err).Error("failed to fetch music")
		return nil, fmt.Errorf("failed to send http: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("breed not found")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.WithError(err).Error("failed to read body")
		return nil, err
	}
	log.Debug("body:", string(body))
	if string(body) == "INVALID_DATA" {
		return nil, fmt.Errorf("breed not found")
	}
	if err = json.Unmarshal(body, &cat); err != nil {
		// log.WithError(err).Error("failed to unmarshal")
		return nil, fmt.Errorf("failed to unmarshal %v into %+v: ", string(body), cat)
	}
	return &cat, nil
}
