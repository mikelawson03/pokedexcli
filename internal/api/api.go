package api

import (
	"fmt"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"github.com/mikelawson03/pokedexcli/internal/pokecache"
	"github.com/mikelawson03/pokedexcli/internal/pokedex"
)

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Count    int         `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []Location  `json:"results"`
}

type Client struct {
	nextURL     *string
	previousURL *string
	cache       *pokecache.Cache
}

type EncountersResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}



func (c *Client) GetNextLocations() (LocationArea, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	
	if c.nextURL != nil {
		url = *c.nextURL
	}

	body, err := c.fetch(url)
	if err != nil {
		return LocationArea{}, err
	}
	locations, err := c.unmarshalLocation(body)

	return locations, err
}

func (c *Client) GetPreviousLocations() (LocationArea, error) {	
	if c.previousURL == nil {
		return LocationArea{}, fmt.Errorf("you're on the first page")
	}

	body, err := c.fetch(*c.previousURL)
	if err != nil {
		return LocationArea{}, err
	}
	locations, err := c.unmarshalLocation(body)

	return locations, err
}

func (c *Client) unmarshalLocation(body []byte) (LocationArea, error) {
	var locations LocationArea
	if err := json.Unmarshal(body, &locations); err != nil {
		return LocationArea{}, fmt.Errorf("error unmarshaling data: %w", err)
	}

	c.nextURL = locations.Next
	c.previousURL = locations.Previous

	return locations, nil
}

func (c *Client) unmarshalEncounters(body []byte) (EncountersResponse, error) {
	var encounters EncountersResponse
	if err := json.Unmarshal(body, &encounters); err != nil {
		return EncountersResponse{}, fmt. Errorf("error unmarshaling data: %w", err)
	}

	return encounters, nil
}

func (c *Client) unmarshalPokemon(body []byte) (pokedex.Pokemon, error) {
	var pokemon pokedex.Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokedex.Pokemon{}, fmt.Errorf("error unmarshaling data: %w", err)
	}

	return pokemon, nil
}

func (c *Client) GetEncounters(location string) (EncountersResponse, error) {

	url := "https://pokeapi.co/api/v2/location-area/" + location + "/"

	body, err := c.fetch(url)
	if err != nil {
		return EncountersResponse{}, fmt.Errorf("error fetching data: %w", err)
	}

	encounters, err := c.unmarshalEncounters(body)

	return encounters, err
	
}

func (c *Client) GetPokemon(name string) (pokedex.Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + name + "/"
	body, err := c.fetch(url)
	if err != nil {
		return pokedex.Pokemon{}, err
	}

	pokemon, err := c.unmarshalPokemon(body)

	return pokemon, err
}

func (c *Client) fetch(url string) ([]byte, error) {
	var body [] byte
	cached, ok := c.cache.Get(url)
	if ok{
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer res.Body.Close()
		
		body, err = io.ReadAll(res.Body); 
		if err != nil {
			return nil, fmt.Errorf("error reading HTTP response data: %w", err)
		}

		c.cache.Add(url, body)
	}
	return body, nil
}

func NewClient(interval time.Duration) *Client {
	c := &Client{
		cache: pokecache.NewCache(interval * time.Second),
	}

	return c
}

