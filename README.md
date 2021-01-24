## Experiments in machine learning for Mahjong

#### Directory Overview

- `bin` precompiled binaries
- `server` code for the mahjong server
- `api` python api code to interact with the server
- `ml` files related to rl algorithm

### Server

Starts on localhost port `8000`, change this by setting the `PORT` env variable.

#### Server API
- `GET /` Server index
- `GET /new` Create a new game, returns the game id and the location.
- `GET /game/<id>` View the (human readable) game state
- `POST /game/<id>` Update the game state. Requires POST data to contain a map with as keys the players (1-indexed) required to perform an action in the current state and as values the index of the action to be performed by that player. 

### API

Tested with python 3.7.