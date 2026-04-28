package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AllLocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type LocationAreaResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonResponse struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

const baseUrl string = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache      *Cache
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: 5 * time.Second,
		},
		cache: NewCache(5 * time.Second),
	}
}

func (c *Client) GetAllLocationAreas(pageUrl *string) (AllLocationAreasResponse, error) {
	apiUrl := baseUrl + "/location-area"
	if pageUrl != nil {
		apiUrl = *pageUrl
	}

	return fetch[AllLocationAreasResponse](c, apiUrl)
}

func (c *Client) GetLocationArea(name string) (LocationAreaResponse, error) {
	return fetch[LocationAreaResponse](c, baseUrl+"/location-area/"+name)
}

func (c *Client) GetPokemon(name string) (PokemonResponse, error) {
	return fetch[PokemonResponse](c, baseUrl+"/pokemon/"+name)
}

func fetch[T any](client *Client, url string) (T, error) {
	var result T

	data, ok := client.cache.Get(url)
	if !ok {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return result, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := client.httpClient.Do(req)
		if err != nil {
			return result, fmt.Errorf("failed to fetch: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return result, fmt.Errorf("api responded with status code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return result, fmt.Errorf("failed to read response: %w", err)
		}
		client.cache.Add(url, data)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}
