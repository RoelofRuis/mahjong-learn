from urllib import request, parse

SERVER_PORT = 8000
SERVER_URL = f"http://localhost:{SERVER_PORT}"


def new_game():
    request.urlopen(f"{SERVER_URL}/new")


def send_data(game_id: int, data: map):
    encoded_data = parse.urlencode(data).encode()
    req = request.Request(f"{SERVER_URL}/game/{game_id}", method="POST", data=encoded_data)
    request.urlopen(req)


new_game()
send_data(1, {"1": 0})
