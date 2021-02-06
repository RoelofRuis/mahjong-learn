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

**- `GET /game/<id>` View the (human readable) game state.**
  
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
**- `POST /game/<id>` Update the game state.** Requires POST data to contain a map with as keys the players (1-indexed) required to perform an action in the current state and as values the index of the action to be performed by that player. 

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

**- `GET /game/<id>/player/<player>` View the (human readable) player state**. This contains only part of the state that is visible to the selected player.
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

### API

Tested with python 3.7.