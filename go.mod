module github.com/am1macdonald/pokedexcli

go 1.22.1

replace github.com/am1macdonald/apiLink v0.0.0 => ./internal/apiLink
replace github.com/am1macdonald/locationArea v0.0.0 => ./internal/locationArea
replace github.com/am1macdonald/pokecache v0.0.0 => ./internal/pokecache

require (
	github.com/am1macdonald/apiLink v0.0.0
	github.com/am1macdonald/locationArea v0.0.0
	github.com/am1macdonald/pokecache v0.0.0
)
