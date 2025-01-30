# Pokedex

A command-line Pokedex using the [PokeAPI](https://pokeapi.co/). Explore the Pokemon world, catch Pokemon, and build your collection!

## How to Use

Start the program:
```bash
go run .
```

Available commands:
- `help` - Show all available commands
- `exit` - Close the Pokedex
- `map` - Show the next 20 areas (use multiple times to explore further)
- `mapb` - Go back to previous 20 areas
- `explore <area>` - List all Pokemon in a specific area
- `catch <pokemon>` - Try to catch a Pokemon (harder to catch stronger Pokemon!)
- `inspect <pokemon>` - See details of a Pokemon you've caught
- `pokedex` - See all Pokemon you've caught

Example session:
```
Pokedex > map
canalave-city-area
eterna-city-area
...

Pokedex > explore canalave-city-area
Found Pokemon:
- tentacool
- tentacruel
...

Pokedex > catch tentacool
Throwing a Pokeball at tentacool...
tentacool was caught!

Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
Stats:
  -hp: 40
  -attack: 40
  -defense: 35
Types:
  - water
  - poison
```

## Technical Notes

- Uses caching to remember responses for 5 minutes
- Catch rates are based on Pokemon's base experience
- All data comes from PokeAPI - no local storage
- Commands are case-insensitive
