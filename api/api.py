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


def new_game() -> dict:
    return _do_request(f"{SERVER_URL}/new")


def get_state(game_id: int) -> dict:
    return _do_request(f"{SERVER_URL}/game/{game_id}")


def send_actions(game_id: int, data: map):
    encoded_data = parse.urlencode(data).encode()
    req = request.Request(f"{SERVER_URL}/game/{game_id}", method="POST", data=encoded_data)
    request.urlopen(req)
