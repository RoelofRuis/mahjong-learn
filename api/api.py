from urllib import request, parse
import json

SERVER_PORT = 8000
SERVER_URL = f"http://localhost:{SERVER_PORT}"


def _do_request(url: str) -> dict:
    resp = request.urlopen(url)
    if resp.getcode() == 200 or resp.getcode() == 201:
        data = resp.read()
        return json.loads(data)
    else:
        raise Exception("received non 200 response")


def new_game():
    return _do_request(f"{SERVER_URL}/new")


def get_state(game_id: int):
    return _do_request(f"{SERVER_URL}/game/{game_id}")


def send_data(game_id: int, data: map):
    encoded_data = parse.urlencode(data).encode()
    req = request.Request(f"{SERVER_URL}/game/{game_id}", method="POST", data=encoded_data)
    request.urlopen(req)


# start new game
id = new_game()["id"]

game_state = get_state(id)
print(game_state)

# pick the first action for player 1
send_data(id, {"1": 0})

# pick the first actions for all players
send_data(id, {"2": 0, "3": 0, "4": 0})
