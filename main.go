package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Galaxybigbird/Pokedex/internal/pokecache"
)

type config struct {
	nextLocationAreaURL     *string
	previousLocationAreaURL *string
	cache                   *pokecache.Cache
	lastArgs                []string
	caughtPokemon           map[string]Pokemon
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
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

type locationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.nextLocationAreaURL != nil {
		url = *cfg.nextLocationAreaURL
	}

	// Check cache first
	if data, ok := cfg.cache.Get(url); ok {
		fmt.Println("Cache hit!")
		var result locationAreasResp
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}

		cfg.nextLocationAreaURL = result.Next
		cfg.previousLocationAreaURL = result.Previous

		for _, area := range result.Results {
			fmt.Println(area.Name)
		}
		return nil
	}

	fmt.Println("Cache miss!")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Store in cache
	cfg.cache.Add(url, body)

	var result locationAreasResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	cfg.nextLocationAreaURL = result.Next
	cfg.previousLocationAreaURL = result.Previous

	for _, area := range result.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.previousLocationAreaURL == nil {
		fmt.Println("You're on the first page!")
		return nil
	}

	url := *cfg.previousLocationAreaURL

	// Check cache first
	if data, ok := cfg.cache.Get(url); ok {
		fmt.Println("Cache hit!")
		var result locationAreasResp
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}

		cfg.nextLocationAreaURL = result.Next
		cfg.previousLocationAreaURL = result.Previous

		for _, area := range result.Results {
			fmt.Println(area.Name)
		}
		return nil
	}

	fmt.Println("Cache miss!")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Store in cache
	cfg.cache.Add(url, body)

	var result locationAreasResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	cfg.nextLocationAreaURL = result.Next
	cfg.previousLocationAreaURL = result.Previous

	for _, area := range result.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location area name")
	}
	locationAreaName := args[0]

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationAreaName)

	// Check cache first
	if data, ok := cfg.cache.Get(url); ok {
		fmt.Println("Cache hit!")
		var result locationArea
		err := json.Unmarshal(data, &result)
		if err != nil {
			return err
		}

		fmt.Printf("Exploring %s...\n", locationAreaName)
		fmt.Println("Found Pokemon:")
		for _, encounter := range result.PokemonEncounters {
			fmt.Printf("- %s\n", encounter.Pokemon.Name)
		}
		return nil
	}

	fmt.Println("Cache miss!")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Store in cache
	cfg.cache.Add(url, body)

	var result locationArea
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationAreaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range result.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	// Check cache first
	var pokemon Pokemon
	if data, ok := cfg.cache.Get(url); ok {
		fmt.Println("Cache hit!")
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Cache miss!")
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(body, &pokemon)
		if err != nil {
			return err
		}

		// Store in cache
		cfg.cache.Add(url, body)
	}

	// Calculate catch probability - higher base experience means lower chance
	catchRate := 0.7 - float64(pokemon.BaseExperience)/1000.0
	if catchRate < 0.1 {
		catchRate = 0.1
	}

	// Random chance to catch
	r := rand.Float64()
	if r > catchRate {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	// Pokemon was caught!
	fmt.Printf("%s was caught!\n", pokemonName)
	cfg.caughtPokemon[pokemonName] = pokemon
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}
	pokemonName := args[0]

	pokemon, ok := cfg.caughtPokemon[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("pokedex command takes no arguments")
	}

	fmt.Println("Your Pokedex:")
	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("Empty!")
		return nil
	}

	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lists the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokemon in a location area",
			callback: func(cfg *config) error {
				return commandExplore(cfg, cfg.lastArgs)
			},
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback: func(cfg *config) error {
				return commandCatch(cfg, cfg.lastArgs)
			},
		},
		"inspect": {
			name:        "inspect",
			description: "View information about a caught pokemon",
			callback: func(cfg *config) error {
				return commandInspect(cfg, cfg.lastArgs)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the Pokemon you have caught",
			callback: func(cfg *config) error {
				return commandPokedex(cfg, cfg.lastArgs)
			},
		},
	}
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	commands := getCommands()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	// Convert to lowercase and trim spaces
	text = strings.ToLower(strings.TrimSpace(text))
	// Split on whitespace
	words := strings.Fields(text)
	return words
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		cache:         pokecache.NewCache(5 * time.Minute),
		caughtPokemon: make(map[string]Pokemon),
	}
	rand.Seed(time.Now().UnixNano())
	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := words[1:]

		command, exists := commands[commandName]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		cfg.lastArgs = args
		err := command.callback(cfg)
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
		}
	}
}
