package pokemonapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type PokemonAPI struct {
	cache pokecache.Cache
}

func NewApi() *PokemonAPI {
	return &PokemonAPI{cache: *pokecache.NewCache(5 * time.Second)}
}

func (api *PokemonAPI) GetLocation(url string) (locations []string, next string, prev string, err error) {
	v, isPresent := api.cache.Get(url)

	var apiResponse LocationAreaResponse

	if isPresent {
		fmt.Printf("Cache hit for %v\n", url)
		if err := json.Unmarshal(v, &apiResponse); err != nil {
			return []string{}, "", "", err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, "", "", err
		}

		if res.StatusCode != http.StatusOK {
			return []string{}, "", "", fmt.Errorf("Error: Got %v", res.Status)
		}

		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)

		api.cache.Add(url, bodyBytes)

		if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
			return []string{}, "", "", err
		}
	}

	locations = []string{}
	for _, location := range apiResponse.Results {
		locations = append(locations, location.Name)
	}

	return locations, apiResponse.Next, apiResponse.Previous, nil
}

func (api *PokemonAPI) GetLocationAreaDetail(url string) (pokemons []string, err error) {
	v, isPresent := api.cache.Get(url)

	var apiResponse LocationAreaDetailResponse

	if isPresent {
		fmt.Printf("Cache hit for %v\n", url)
		if err := json.Unmarshal(v, &apiResponse); err != nil {
			return []string{}, err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return []string{}, err
		}

		if res.StatusCode != http.StatusOK {
			return []string{}, fmt.Errorf("Error: Got %v", res.Status)
		}

		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)

		api.cache.Add(url, bodyBytes)

		if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
			return []string{}, err
		}
	}

	pokemons = []string{}
	for _, encounter := range apiResponse.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}

	return pokemons, nil
}
