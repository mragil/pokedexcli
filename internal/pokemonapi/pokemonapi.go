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

func GetRequest[T any](api *PokemonAPI, url string) (T, error) {
	v, isPresent := api.cache.Get(url)

	var apiResponse T

	if isPresent {
		fmt.Printf("Cache hit for %v\n", url)
		if err := json.Unmarshal(v, &apiResponse); err != nil {
			return apiResponse, err
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return apiResponse, err
		}

		if res.StatusCode != http.StatusOK {
			return apiResponse, fmt.Errorf("Error: Got %v", res.Status)
		}

		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)

		api.cache.Add(url, bodyBytes)

		if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
			return apiResponse, err
		}
	}

	return apiResponse, nil
}

func (api *PokemonAPI) GetLocation(url string) (locations []string, next string, prev string, err error) {
	apiResponse, err := GetRequest[LocationAreaResponse](api, url)

	if err != nil {
		return []string{}, "", "", err
	}

	locations = []string{}
	for _, location := range apiResponse.Results {
		locations = append(locations, location.Name)
	}

	return locations, apiResponse.Next, apiResponse.Previous, nil
}

func (api *PokemonAPI) GetLocationAreaDetail(url string) (pokemons []string, err error) {
	apiResponse, err := GetRequest[LocationAreaDetailResponse](api, url)

	if err != nil {
		return []string{}, err
	}

	pokemons = []string{}
	for _, encounter := range apiResponse.PokemonEncounters {
		pokemons = append(pokemons, encounter.Pokemon.Name)
	}

	return pokemons, nil
}

func (api *PokemonAPI) GetPokemonDetail(url string) (pokemon Pokemon, err error) {
	apiResponse, err := GetRequest[PokemonResponse](api, url)

	if err != nil {
		return Pokemon{}, err
	}

	types := []string{}
	for _, v := range apiResponse.Types {
		types = append(types, v.Type.Name)
	}

	stats := PokemonStats{}
	statsMap := map[string]*int{
		"hp":              &stats.HP,
		"attack":          &stats.Attack,
		"defense":         &stats.Defense,
		"speed":           &stats.Speed,
		"special-attack":  &stats.SpecialAttack,
		"special-defense": &stats.SpecialDefense,
	}

	for _, v := range apiResponse.Stats {
		if fieldPtr, ok := statsMap[v.Stat.Name]; ok {
			*fieldPtr = v.BaseStat
		}
	}

	return Pokemon{
		Name:   apiResponse.Name,
		Height: apiResponse.Height,
		Weight: apiResponse.Weight,
		Types:  types,
		Stats:  stats,
		Chance: apiResponse.BaseExperience,
	}, nil
}
