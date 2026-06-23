package pokemonapi

type BaseResultResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaDetailResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod BaseResultResponse `json:"encounter_method"`
		VersionDetails  []struct {
			Rate    int                `json:"rate"`
			Version BaseResultResponse `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int                `json:"game_index"`
	ID        int                `json:"id"`
	Location  BaseResultResponse `json:"location"`
	Name      string             `json:"name"`
	Names     []struct {
		Language BaseResultResponse `json:"language"`
		Name     string             `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon        BaseResultResponse `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int                `json:"chance"`
				ConditionValues []any              `json:"condition_values"`
				MaxLevel        int                `json:"max_level"`
				Method          BaseResultResponse `json:"method"`
				MinLevel        int                `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int                `json:"max_chance"`
			Version   BaseResultResponse `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type LocationAreaResponse struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []BaseResultResponse `json:"results"`
}

type PokemonResponse struct {
	ID             int                    `json:"id"`
	Name           string                 `json:"name"`
	Height         int                    `json:"height"`
	Weight         int                    `json:"weight"`
	Stats          []PokemonStatsResponse `json:"stats"`
	BaseExperience int                    `json:"base_experience"`
	Types          []PokemonTypesResponse `json:"types"`
}

type PokemonStatsResponse struct {
	BaseStat int                `json:"base_stat"`
	Effort   int                `json:"effort"`
	Stat     BaseResultResponse `json:"stat"`
}

type PokemonTypesResponse struct {
	Slot int                `json:"slot"`
	Type BaseResultResponse `json:"type"`
}

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  PokemonStats
	Chance int
	Types  []string
}
type PokemonStats struct {
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
}

type PokemonTrainer struct {
	CatchChance int
	Pokedex     map[string]*Pokemon
}
