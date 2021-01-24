from urllib import request, parse

SERVER_PORT = 8000
SERVER_URL = f"http://localhost:{SERVER_PORT}"


def new_game():
    request.urlopen(f"{SERVER_URL}/new")


def send_data(game_id: int, data: map):
    encoded_data = parse.urlencode(data).encode()
    req = request.Request(f"{SERVER_URL}/game/{game_id}", method="POST", data=encoded_data)
    request.urlopen(req)


# start new game
new_game()

# pick the first action for player 1
send_data(1, {"1": 0})

# pick the first actions for all players
send_data(1, {"1": 0, "2": 0, "3": 0, "4": 0})
