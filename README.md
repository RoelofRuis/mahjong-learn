## Experiments in machine learning for Mahjong

### Directory Overview

- `bin` precompiled binaries
- `server` code for the mahjong server
- `api` python api code to interact with the server
- `ml` files related to rl algorithm

### Server

Starts on localhost port `8000`, change this by setting the `PORT` env variable.

#### Server API

All API requests return JSON

**`GET /` Server index.**

```
Status Code 200
{
    message:       string
    version:       string
    games_started: int
    new_game:      url
}
```

**- `GET /new` Create a new game, returns the game id and the location.**

```
Status Code 201
{
    message:  string
    id:       int
    location: url    
}
```

**- `GET /game/<id>` View the human readable game state.**

```
Status Code 200
{
    has_ended:      bool
    state_name:     string
    prevalent_wind: string
    active_players: []int
    active_discard: string
    players:        string -> {
        actions:   string -> string
        score:     int
        wind:      string
        received:  string
        concealed: []string
        exposed:   []string
        discarded: []string    
    }
    wall:           []string        
}
```

**- `POST /game/<id>` Update the game state.** Requires POST data to contain a map with as keys the players (0-indexed)
required to perform an action in the current state and as values the index of the action to be performed by that player.

```
Status Code 202
{
    message:  string
    id:       int
    location: url
}

Status Code 400 (In case an incorrect action was sent)
{
    error:       string
    status_code: int
}
```

**- `GET /game/<id>/player/<player>` View the human readable player state**. This contains only part of the state that
is visible to the selected player.

```
Status Code 200
{
    actions:           string -> string
    prevalent_wind:    string
    discarding_player: int
    active_discard:    string
    score:             int
    wind:              string
    received:          string
    concealed:         []string
    exposed:           []string
    discarded:         []string
    other_players:     string -> {
        score:     int
        wind:      string
        exposed:   []string
        discarded: []string    
    }        
}
```

**- `GET /game/<id>/player/<player>?vec=1` View the vectorized player state**.

```
Status Code 200
{
    score:                         int       1
    bonus_tiles:                   []int     1x8
    prevalent_wind:                []int     1x4
    player_wind:                   []int     1x4
    discarding_player:             []int     1x3
    active_discard:                []int     1x3   [tile]
    received_tile:                 []int     1x3   [tile]
    concealed_tiles:               [][]int   13x3  [tile]
    exposed_chows:                 [][][]int 4x3x3 [tile]
    exposed_pungs:                 [][]int   4x3   [tile]
    exposed_kongs:                 [][]int   4x3   [tile]
    hidden_kongs:                  [][]int   4x3   [tile]
    discards:                      [][]int   40x3  [tile]
    right_player_score:            int       1
    right_player_bonus_tiles:      []int     1x8
    right_player_wind:             []int     1x4
    right_player_exposed_chows:    [][][]int 4x3x3 [tile]
    right_player_exposed_pungs:    [][]int   4x3   [tile]
    right_player_exposed_kongs:    [][]int   4x3   [tile]
    right_player_discards:         [][]int   40x3  [tile]
    opposite_player_score:         int       1
    opposite_player_bonus_tiles:   []int     1x8
    opposite_player_wind:          []int     1x4
    opposite_player_exposed_chows: [][][]int 4x3x3 [tile]
    opposite_player_exposed_pungs: [][]int   4x3   [tile]
    opposite_player_exposed_kongs: [][]int   4x3   [tile]
    opposite_player_hidden_kongs:  [][]int   4x3   [tile]
    opposite_player_discards:      [][]int   40x3  [tile]
    left_player_score:             int       1
    left_player_bonus_tiles:       []int     1x8
    left_player_wind:              []int     1x4
    left_player_exposed_chows:     [][][]int 4x3x3 [tile]
    left_player_exposed_pungs:     [][]int   4x3   [tile]
    left_player_exposed_kongs:     [][]int   4x3   [tile]
    left_player_hidden_kongs:      [][]int   4x3   [tile]
    left_player_discards:          [][]int   40x3  [tile]
}
```

### API

Tested with python 3.7.