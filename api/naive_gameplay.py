import api

# start new game
id = api.new_game()["id"]

while True:
    game_state = api.get_state(id)
    state_name = game_state["state_name"]
    active_player = game_state["active_player"]
    print(f"{state_name}: {active_player}")

    # pick the first action for player 1
    api.send_actions(id, {f"{active_player}": 0})

    # pick the first actions for other players
    api.send_actions(id, {f"{p}": 0 for p in [1, 2, 3, 4] if p != active_player})
