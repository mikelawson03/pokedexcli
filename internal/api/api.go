package api

import (
	"encoding/json"
	"fmt"
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
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
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

type PokemonResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
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
		return EncountersResponse{}, fmt.Errorf("error unmarshaling data: %w", err)
	}

	return encounters, nil
}

func (c *Client) unmarshalPokemon(body []byte) (PokemonResponse, error) {
	var pokemon PokemonResponse
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return PokemonResponse{}, fmt.Errorf("error unmarshaling data: %w", err)
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

	resp, err := c.unmarshalPokemon(body)
	if err != nil {
		return pokedex.Pokemon{}, err
	}

	pokemon := transformPokemon(resp)
	return pokemon, err
}

func transformPokemon(resp PokemonResponse) pokedex.Pokemon {
	pokemon := pokedex.Pokemon{
		Name:           resp.Name,
		BaseExperience: resp.BaseExperience,
		Height:         resp.Height,
		Weight:         resp.Weight,
		Stats:          make(map[string]int),
		Types:          []string{},
	}

	for _, v := range resp.Stats {
		pokemon.Stats[v.Stat.Name] = v.BaseStat
	}

	for _, v := range resp.Types {
		pokemon.Types = append(pokemon.Types, v.Type.Name)
	}

	return pokemon
}

func (c *Client) fetch(url string) ([]byte, error) {
	var body []byte
	cached, ok := c.cache.Get(url)
	if ok {
		body = cached
	} else {
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making request: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
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
