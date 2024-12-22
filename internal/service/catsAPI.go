package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewCatsApiService(client *http.Client) *CatsApiService {
	return &CatsApiService{
		c: client,
	}
}

type CatsApiService struct {
	c *http.Client
}

type CatPicture struct {
	Breeds []any  `json:"breeds"`
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CatApiResponse struct {
	Weight struct {
		Imperial string `json:"imperial"`
		Metric   string `json:"metric"`
	} `json:"weight"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	CfaURL           string `json:"cfa_url,omitempty"`
	VetstreetURL     string `json:"vetstreet_url,omitempty"`
	VcahospitalsURL  string `json:"vcahospitals_url,omitempty"`
	Temperament      string `json:"temperament"`
	Origin           string `json:"origin"`
	CountryCodes     string `json:"country_codes"`
	CountryCode      string `json:"country_code"`
	Description      string `json:"description"`
	LifeSpan         string `json:"life_span"`
	Indoor           int    `json:"indoor"`
	Lap              int    `json:"lap,omitempty"`
	AltNames         string `json:"alt_names,omitempty"`
	Adaptability     int    `json:"adaptability"`
	AffectionLevel   int    `json:"affection_level"`
	ChildFriendly    int    `json:"child_friendly"`
	DogFriendly      int    `json:"dog_friendly"`
	EnergyLevel      int    `json:"energy_level"`
	Grooming         int    `json:"grooming"`
	HealthIssues     int    `json:"health_issues"`
	Intelligence     int    `json:"intelligence"`
	SheddingLevel    int    `json:"shedding_level"`
	SocialNeeds      int    `json:"social_needs"`
	StrangerFriendly int    `json:"stranger_friendly"`
	Vocalisation     int    `json:"vocalisation"`
	Experimental     int    `json:"experimental"`
	Hairless         int    `json:"hairless"`
	Natural          int    `json:"natural"`
	Rare             int    `json:"rare"`
	Rex              int    `json:"rex"`
	SuppressedTail   int    `json:"suppressed_tail"`
	ShortLegs        int    `json:"short_legs"`
	WikipediaURL     string `json:"wikipedia_url,omitempty"`
	Hypoallergenic   int    `json:"hypoallergenic"`
	ReferenceImageID string `json:"reference_image_id,omitempty"`
	CatFriendly      int    `json:"cat_friendly,omitempty"`
	Bidability       int    `json:"bidability,omitempty"`
}

func (s *CatsApiService) GetRandomCatPicture(ctx context.Context) ([]byte, error) {
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
	resp.Body.Close()

	resp, err = s.c.Get(cats[0].URL)
	if err != nil {
		// log.WithError(err).Error("failed to fetch music")
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		// log.WithError(err).Error("failed to read body")
		return nil, err
	}

	return body, nil
}

func (s *CatsApiService) GetAllBreeds(ctx context.Context) ([]string, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.GetAllBreeds")

	var cats []CatApiResponse
	fmt.Println(s.c.Timeout)
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

	var breeds []string
	for _, cat := range cats {
		breeds = append(breeds, cat.Name+" - "+cat.ID)
	}

	if len(breeds) == 0 {
		// log.Error("no cats found")
		return nil, fmt.Errorf("no cats found")

	}
	return breeds, nil
}

func (s *CatsApiService) GetBreed(ctx context.Context, breed string) (*CatApiResponse, error) {
	// log := logger.GetLoggerFromCtx(ctx).WithField("op", "UserService.GetAllBreeds")

	var cat CatApiResponse
	resp, err := s.c.Get("https://api.thecatapi.com/v1/breeds/" + breed)
	if err != nil {
		// log.WithError(err).Error("failed to fetch music")
		return nil, fmt.Errorf("failed to send http: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.WithError(err).Error("failed to read body")
		return nil, err
	}
	if string(body) == "INVALID_DATA" {
		return nil, fmt.Errorf("breed not found")
	}
	if err = json.Unmarshal(body, &cat); err != nil {
		// log.WithError(err).Error("failed to unmarshal")
		return nil, fmt.Errorf("failed to unmarshal %v into %+v: ", string(body), cat)
	}
	return &cat, nil
}
