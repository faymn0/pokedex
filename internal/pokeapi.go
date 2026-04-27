package pokeapi 

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseUrl string = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Duration(5) * time.Second,
		},
	}
}

func (c *Client) GetLocationArea(pageUrl *string) (locationAreaResponse, error) {
	apiUrl := baseUrl + "/location-area"
	if pageUrl != nil {
		apiUrl = *pageUrl
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req) 
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to fetch: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return locationAreaResponse{}, fmt.Errorf("api responded with status code: %d", res.StatusCode)
	}

	var locationAreas locationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locationAreas); err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return locationAreas, nil
} 

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}
