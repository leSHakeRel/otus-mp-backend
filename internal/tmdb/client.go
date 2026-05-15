package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"movie-night-planner-backend/internal/config"
)

type Client struct {
	apiKey       string
	baseURL      string
	imageBaseURL string
	httpClient   *http.Client
}

type Movie struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	PosterPath   string     `json:"poster_path"`
	BackdropPath string     `json:"backdrop_path"`
	ReleaseDate  *time.Time `json:"release_date"`
	VoteAverage  float64    `json:"vote_average"`
	Overview     string     `json:"overview"`
}

type SearchResponse struct {
	Page         int     `json:"page"`
	TotalResults int     `json:"total_results"`
	TotalPages   int     `json:"total_pages"`
	Results      []Movie `json:"results"`
}

type MovieDetails struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	PosterPath   string     `json:"poster_path"`
	BackdropPath string     `json:"backdrop_path"`
	ReleaseDate  *time.Time `json:"release_date"`
	VoteAverage  float64    `json:"vote_average"`
	Overview     string     `json:"overview"`
}

func NewClient(cfg *config.TMDBConfig) *Client {
	return &Client{
		apiKey:       cfg.APIKey,
		baseURL:      cfg.BaseURL,
		imageBaseURL: cfg.ImageBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SearchMovies(query string, page int) (*SearchResponse, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&language=en-US&page=%d&query=%s",
		c.baseURL, c.apiKey, page, query)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb api error: %s - %s", resp.Status, string(body))
	}

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	return &searchResp, nil
}

func (c *Client) GetMovieDetails(tmdbID int) (*MovieDetails, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s&language=en-US",
		c.baseURL, tmdbID, c.apiKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tmdb api error: %s - %s", resp.Status, string(body))
	}

	var movie MovieDetails
	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (c *Client) GetImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	if size == "" {
		size = "w500"
	}
	return fmt.Sprintf("%s/%s%s", c.imageBaseURL, size, path)
}
